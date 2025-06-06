---
description: Learn how to work effectively with Warehouses
sidebar_label: Working with Warehouses
---

# Working with Warehouses

Kargo `Warehouse` resources each manage subscriptions to one or more of various
types of artifact sources, including:

- Container image repositories
- Git repositories
- Helm chart repositories

When a `Warehouse` observes a new revision of any artifact to which it
subscribes, it creates a new `Freight` resource representing a specific
collection of artifact revisions that can be promoted from `Stage` to `Stage`
_as a unit_.

:::info
For a broader, conceptual understanding of warehouses and their relation
to other Kargo concepts, refer to
[Core Concepts](./../10-core-concepts/index.md).
:::

## The `Warehouse` Resource Type

A `Warehouse`'s subscriptions are all defined within its `spec.subscriptions`
field.

In this example, a `Warehouse` subscribes to both a container image repository
and a Git repository:

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: my-warehouse
  namespace: kargo-demo
spec:
  subscriptions:
  - image:
      repoURL: public.ecr.aws/nginx/nginx
      semverConstraint: ^1.26.0
  - git:
      repoURL: https://github.com/example/kargo-demo.git
```

The remainder of this document focuses on the configuration of the individual
subscription types.

### Container Image Subscriptions

Container image repository subscriptions can be defined using the following
fields:

- `repoURL`: The URL of the container image repository _without any tag_. This field is required.

- `imageSelectionStrategy`: One of four pre-defined strategies for selecting the
  desired image. (See next section.)

- `allowTags`: An optional regular expression that limits the eligibility for
  selection to tags that match the pattern.

- `ignoreTags`: An optional list of tags that should explicitly be ignored.

- `platform`: An optional identifier that constrains image selection to those
  images supporting the specified operating system and system architecture.
  e.g., `linux/amd64`.

    :::note
    It is seldom necessary to specify this field.
    :::

- `discoveryLimit`: Many selection strategies (see next section) do not actually
  select a _single_ image; rather they select the n best fits for the specified
  constraints. The _best_ fit is the zero element in the list of selected
  images. `discoveryLimit` specifies how many images to discover.
  
    The default is `20`.

    :::note
    For poorly performing `Warehouse`s -- for instance ones frequently
    encountering rate limits -- decreasing this limit may improve performance.
    :::

- `gitRepoURL`: (Deprecated as of version 1.5, will be removed in version 1.7) 
  Optional metadata to inform Kargo of the Git repository that contains the 
  image's Dockerfile or other build context.

- `insecureSkipTLSVerify`: Set to `true` to disable validation of the
  repository's TLS certificate.

    :::warning
    This is a security risk and should only be used in development environments.
    :::

#### Image Selection Strategies

For subscriptions to container image repositories, the `imageSelectionStrategy`
field specifies a method for selecting the desired image. The available
strategies are:

- `SemVer`: Selects the image with the tag that best matches a semantic
  versioning constraint specified by the `semverConstraint` field. With no
  constraint specified, the strategy simply selects the image with the
  semantically greatest tag. All tags that are not valid semantic versions are
  ignored.

   The `strictSemvers` field defaults to `true`, meaning only tags containing
   all three parts of a semantic version (major, minor, and patch) are
   considered. Disabling this should be approached with caution because any
   image tagged only with decimal characters will be considered a valid semantic
   version (containing only the major element).

    __`SemVer` is the default strategy if one is not specified.__

    :::info
    Kargo uses the [semver](https://github.com/masterminds/semver) package for
    parsing and comparing semantic versions and semantic version constraints.
    Refer to
    [these docs](https://github.com/masterminds/semver#checking-version-constraints)
    for detailed information on version constraint syntax.
    :::

    Example:

    ```yaml
    spec:
      subscriptions:
      - image:
          repoURL: public.ecr.aws/nginx/nginx
          semverConstraint: ^1.26.0
    ```

- `Lexical`: This strategy selects the image with the lexicographically greatest
   tag.

   This is useful in scenarios wherein tags incorporate date/time stamps in
   formats such as `yyyymmdd` and you wish to select the tag with the latest
   stamp. When using this strategy, it is recommended to use the `allowTags`
   field to limit eligibility to tags that match the expected format.

    Example:

    ```yaml
    spec:
      subscriptions:
      - image:
          repoURL: public.ecr.aws/nginx/nginx
          imageSelectionStrategy: Lexical
          allowTags: ^nightly-\d{8}$
    ```

- `Digest`: This selects the image _currently_ referenced by some "mutable tag,"
   such as `latest`.
   
    __Unintuitively, the mutable tag name must be specified using the
    `semverConstraint` field.__ Importantly, the _digest_ will change every time
    the tag is updated.

    :::warning
    "Mutable tags": Tags like `latest` that are sometimes, perhaps frequently,
    updated to point to a different, presumably newer image.

    "Immutable tags": Tags that have version or date information embedded within
    them, along with an expectation of never being updated to reference a
    different image.

    Using mutable tags like `latest` _is a widely discouraged practice._
    Whenever possible, it is recommended to use immutable tags.
    :::

    Example:

    ```yaml
    spec:
      subscriptions:
      - image:
          repoURL: public.ecr.aws/nginx/nginx
          imageSelectionStrategy: Digest
          semverConstraint: latest
    ```

- `NewestBuild`: This strategy selects the image with the most recent build
  time.

    :::warning
    `NewestBuild` requires retrieving metadata for every eligible tag, which can
    be slow and is likely to exceed the registry's rate limits. __This can
    result in system-wide performance degradation.__

    If using this strategy is unavoidable, it is recommended to use the
    `allowTags` field to limit the number of tags for which metadata is
    retrieved to reduce the risk of encountering rate limits. `allowTags` may
    require periodic adjustment as a repository grows.
    :::

    ```yaml
    spec:
      subscriptions:
      - image:
          repoURL: public.ecr.aws/nginx/nginx
          imageSelectionStrategy: NewestBuild
          allowTags: ^nightly
    ```

### Git Repository Subscriptions

Git repository subscriptions can be defined using the following fields:

- `repoURL`: The URL of the Git repository. This field is required.

- `commitSelectionStrategy`: One of four pre-defined strategies for selecting
  the desired commit. (See next section.)

- `allowTags`: An optional regular expression that limits the eligibility for
  selection to tags that match the pattern. (This is not applicable to selection
  strategies that do not involve tags.)

- `ignoreTags`: An optional list of tags that should explicitly be ignored.
  (This is not applicable to selection strategies that do not involve tags.)

- `expressionFilter`: An optional expression that filters commits and tags based 
  on their metadata. See [Expression Filtering](#expression-filtering) for 
  details.

- `includePaths`: See
  [Git Subscription Path Filtering](#git-subscription-path-filtering).

- `excludePaths`: See
  [Git Subscription Path Filtering](#git-subscription-path-filtering).

- `discoveryLimit`: Many selection strategies (see next section) do not actually
  select a _single_ commit; rather they select the n best fits for the specified
  constraints. The _best_ fit is the zero element in the list of selected
  commits. `discoveryLimit` specifies how many commits to discover.
  
    The default is `20`.

   :::note
   Lowering this limit for a Git repository subscription does not improve
   performance by the margins that it does for a container image repository
   subscription.
   :::

- `insecureSkipTLSVerify`: Set to `true` to disable validation of the
  repository's TLS certificate.

    :::warning
    This is a security risk and should only be used in development environments.
    :::

#### Commit Selection Strategies

For subscriptions to Git repositories, the `commitSelectionStrategy`
field specifies a method for selecting the desired commit. The available
strategies are:

- `NewestFromBranch`: Selects the most recent commit from a branch specified
  by the `branch` field. If a branch is not specified, the strategy selects
  commits from the repository's default branch (typically `main` or `master`).
  
    This is useful for the average case, wherein you wish for the `Warehouse` to
    continuously discover the latest changes to a branch that receives regular
    updates.

    __`NewestFromBranch` is the default selection strategy if one is not
    specified.__

    Example:

    ```yaml
    spec:
      subscriptions:
      - git:
          repoURL: https://github.com/example/repo.git
          branch: main
    ```

- `SemVer`: Selects the commit referenced by the tag that best matches a
  semantic versioning constraint. All tags that are not valid semantic versions
  are ignored. With no constraint specified, the strategy simply selects the
  commit referenced by the semantically greatest tag.

    This is useful in scenarios wherein you do not wish for the `Warehouse` to
    continuously discover _every new commit_ and would like limit selection to
    commits tagged with a semantic version, and possibly within a certain range.

    The `strictSemvers` field defaults to `true`, meaning only tags containing
    all three parts of a semantic version (major, minor, and patch) are
    considered. Disabling this should be approached with caution because any
    image tagged only with decimal characters will be considered a valid
    semantic version (containing only the major element).

    :::info
    Kargo uses the [semver](https://github.com/masterminds/semver) package for
    parsing and comparing semantic versions and semantic version constraints.
    Refer to
    [these docs](https://github.com/masterminds/semver#checking-version-constraints)
    for detailed information on version constraint syntax.
    :::

    Example:

    ```yaml
    spec:
      subscriptions:
      - git:
          repoURL: https://github.com/example/repo.git
          commitSelectionStrategy: SemVer
          semverConstraint: ^1.0.0
    ```

- `Lexical`: Selects the commit referenced by the lexicographically greatest
  tag.

    This is useful in scenarios wherein you do not wish for the `Warehouse` to
    discover _every new commit_ and tags incorporate date/time stamps in formats
    such as `yyyymmdd` and you wish to select the tag with the latest stamp.
    When using this strategy, it is recommended to use the `allowTags` field to
    limit eligibility to tags that match the expected format.

    Example:

    ```yaml
    spec:
      subscriptions:
      - git:
          repoURL: https://github.com/example/repo.git
          commitSelectionStrategy: Lexical
          allowTags: ^nightly-\d{8}$
    ```

- `NewestTag`: Selects the commit referenced by the most recently created tag.
  
    When using this strategy, it is recommended to use the `allowTags` field to
    limit eligibility to tags that match the expected format.

    Example:

    ```yaml
    spec:
      subscriptions:
      - git:
          repoURL: https://github.com/example/repo.git
          commitSelectionStrategy: NewestTag
          allowTags: ^nightly
    ```

#### Expression Filtering

Git repository subscriptions support advanced filtering using expressions. These
expressions allow you to filter commits and tags based on their metadata using
[expr-lang](https://expr-lang.org) syntax.

:::info
The expressions must evaluate to a boolean value (`true` or `false`). If an
expression evaluates to a non-boolean value, an attempt will be made to
convert it to a boolean (e.g., `0` to `false`, `1` to `true`).
:::

:::warning
Invalid expressions will cause the subscription to fail. Always test your
expressions to ensure they evaluate correctly with your repository's data.
:::

:::tip
You can test your expressions using the
[expr-lang playground](https://expr-lang.org/playground).

The playground allows you to evaluate expressions against sample data and
see the results in real-time. This is especially useful for debugging and
validating your expressions before applying them to your `Warehouse` resources.
:::

The `expressionFilter` field provides a unified way to filter commits or tags
based on the selected commit selection strategy. The behavior and available
variables depend on your `commitSelectionStrategy`:

**For commit-based filtering** (`NewestFromBranch` strategy):

- Filters commits based on commit metadata
- Applied when selecting the newest commit from a branch

**For tag-based filtering** (`SemVer`, `Lexical`, and `NewestTag` strategies):

- Filters tags based on tag and associated commit metadata
- Applied after `allowTags`, `ignoreTags`, and `semverConstraint` fields

##### Available Expression Filtering Variables

The variables available in your expression depend on the commit selection
strategy:

**For `NewestFromBranch` (commit filtering):**

- `id`: The ID (SHA) of the commit
- `commitDate`: The date of the commit
- `author`: The author of the commit, in format `Name <email>`
- `committer`: The committer of the commit, in format `Name <email>`
- `subject`: The first line of the commit message

**For `SemVer`, `Lexical`, and `NewestTag` (tag filtering):**

- `tag`: The name of the tag
- `id`: The commit ID that the tag references
- `creatorDate`: The tag creation date (annotated tag) or commit date
  (lightweight tag)
- `author`: The author of the commit that the tag references, in the format of
  `Name <email>`
- `committer`: The committer of the commit that the tag references, in the
  format of `Name <email>`
- `subject`: The first line of the commit message associated with the tag
- `tagger`: The tagger of the tag, in the format of `Name <email>`. Only
  available for annotated tags.
- `annotation`: The first line of the tag annotation. Only available for
  annotated tags.

##### Expression Filtering Examples

**Filtering commits by excluding bot authors:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: NewestFromBranch
      expressionFilter: !(author contains '<bot@example.com>')
```

**Filtering commits with specific message patterns:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: NewestFromBranch
      expressionFilter: subject contains 'feat:' || subject contains 'fix:'
```

**Filtering commits with multiple criteria:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: NewestFromBranch
      expressionFilter: !(author == 'Example Bot') && commitDate.After(date('2025-01-01'))
```

**Filtering commits to exclude those with ignore markers:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: NewestFromBranch
      expressionFilter: !(subject contains '[kargo-ignore]')
```

**Filtering tags by author name:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: SemVer
      expressionFilter: author == 'John Doe <john@example.com>'
```

**Filtering tags created after a specific date:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: NewestTag
      expressionFilter: creatorDate.Year() >= 2024
```

**Filtering tags to exclude those committed by bots:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: Lexical
      expressionFilter: !(committer contains '<bot@example.com>')
```

**Filtering tags with complex conditions:**

```yaml
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/repo.git
      commitSelectionStrategy: SemVer
      expressionFilter: creatorDate.After(date('2024-01-01')) && !(tag contains 'alpha')
```

#### Git Subscription Path Filtering

In some cases, it may be necessary to constrain the paths within a Git
repository that a `Warehouse` will consider as triggers for `Freight`
production. This is especially useful for GitOps repositories that are
["monorepos"](../30-patterns/index.md#monorepo-layout) containing configuration
for multiple applications.

The paths that may or must not trigger `Freight` production may be specified
using a combination of the `includePaths` and `excludePaths` fields of a Git
repository subscription.

The following example demonstrates a `Warehouse` with a Git repository
subscription that will only produce new `Freight` when the latest commit
(selected by the applicable commit selection strategy) contains changes in the
`apps/guestbook` directory since the last piece of `Freight` produced by the
`Warehouse`:

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: my-warehouse
  namespace: kargo-demo
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/kargo-demo.git
      includePaths:
      - apps/guestbook
```

The next example demonstrates the opposite: a `Warehouse` with a Git repository
subscription that will only produce new `Freight` when the latest commit
(selected by the applicable commit selection strategy) contains changes to paths
_other than_ the repository's `docs/` directory:

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: my-warehouse
  namespace: kargo-demo
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/kargo-demo.git
      excludePaths:
      - docs
```

`includePaths` and `excludePaths` may be combined to include a broad set of
paths and then exclude a subset of those. The following example demonstrates a
`Warehouse` with a Git repository subscription that will only produce new
`Freight` when the latest commit (selected by the applicable commit selection
strategy) contains changes _within_ the `apps/guestbook` directory _other than_
the `apps/guestbook/README.md`:

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: my-warehouse
  namespace: kargo-demo
spec:
  subscriptions:
  - git:
      repoURL: https://github.com/example/kargo-demo.git
      includePaths:
      - apps/guestbook
      excludePaths:
      - apps/guestbook/README.md
```

:::note
It is important to understand that new `Freight` will be produced when the
latest commit (selected by the applicable commit selection strategy) contains
_even a single change_ that is:

1. Implicitly included via undefined `includePaths`.

    OR

    Explicitly included via `includePaths`.

    AND

2. Not explicitly excluded via `excludePaths`.
:::

:::note
By default, the strings in the `includePaths` and `excludePaths` fields are
treated as exact paths to files or directories. (Selecting a directory will
implicitly select all paths within that directory.)

Paths may _also_ be specified using glob patterns (by prefixing the string with
`glob:`) or regular expressions (by prefixing the string with `regex:` or
`regexp:`).
:::

### Helm Chart Repository Subscriptions

Helm chart repository subscriptions can be defined using the following fields:

- `repoURL`: The URL of the Helm chart repository. This field is required.

  Chart repositories using http/s may contain versions of many _different_
  charts. Subscriptions to all chart repositories using http/s __must__
  additionally specify the chart's name in the `name` field.

  For chart repositories in OCI registries, the repository URL points only to
  revisions of a _single_ chart. Subscriptions to chart repositories in OCI
  registries __must__ leave the `name` field empty.

- `name`: See above.

- `semverConstraint`: Selects the chart version best matching this constraint.
  If left unspecified, the subscription implicitly selects the semantically
  greatest version of the chart.

  :::info
  Helm _requires_ charts to be semantically versioned.
  :::

  :::info
  Kargo uses the [semver](https://github.com/masterminds/semver) package for
  parsing and comparing semantic versions and semantic version constraints.
  Refer to
  [these docs](https://github.com/masterminds/semver#checking-version-constraints)
  for detailed information on version constraint syntax.
  :::

- `discoveryLimit`: A chart repository subscription does not actually select a
  _single_ chart version; rather it selects the n best fits for the specified
  constraints. The _best_ fit is the zero element in the list of selected
  charts. `discoveryLimit` specifies how many chart versions to discover.

  The default is `20`.

  Example:

  ```yaml
  spec:
    subscriptions:
    - chart:
        repoURL: https://charts.example.com
        name: my-chart
        semverConstraint: ^1.0.0
  ```
