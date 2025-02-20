// Code generated by quicktype. DO NOT EDIT.

package directives

type CommonDefs interface{}

type ComposeOutput map[string]interface{}

type ArgoCDUpdateConfig struct {
	Apps []ArgoCDAppUpdate `json:"apps"`
}

type ArgoCDAppUpdate struct {
	// Specifies the name of an Argo CD Application resource to be updated.
	Name string `json:"name"`
	// Specifies the namespace of an Argo CD Application resource to be updated. If left
	// unspecified, the namespace will be the controller's configured default.
	Namespace string `json:"namespace,omitempty"`
	// Describes updates to be applied to various sources of an Argo CD Application resource.
	Sources []ArgoCDAppSourceUpdate `json:"sources,omitempty"`
}

type ArgoCDAppSourceUpdate struct {
	// If applicable, identifies a specific chart within the Helm chart repository specified by
	// the 'repoURL' field. When the source to be updated references a Helm chart repository,
	// the values of the 'repoURL' and 'chart' fields should exactly match the values of the
	// same fields in the source. i.e. Do not match the values of these two fields to your
	// Warehouse; match them to the Application source you wish to update.
	Chart string `json:"chart,omitempty"`
	// Specifies the desired revision for the source. If left undefined, the desired revision
	// will be determined by Freight (if possible). Note that the source's 'targetRevision' will
	// not be updated to this commit unless 'updateTargetRevision=true' is set. The utility of
	// this field, on its own, is to specify the revision that the source should be observably
	// synced to during a health check.
	DesiredRevision string                       `json:"desiredRevision,omitempty"`
	Helm            *ArgoCDHelmParameterUpdates  `json:"helm,omitempty"`
	Kustomize       *ArgoCDKustomizeImageUpdates `json:"kustomize,omitempty"`
	// With possible help from the 'chart' field, identifies which of an Argo CD Application's
	// sources is to be updated. When the source to be updated references a Helm chart
	// repository, the values of the 'repoURL' and 'chart' fields should exactly match the
	// values of the same fields in the source. i.e. Do not match the values of these two fields
	// to your Warehouse; match them to the Application source you wish to update.
	RepoURL string `json:"repoURL"`
	// Indicates whether the source should be updated such that its 'targetRevision' field
	// points directly at the desired revision. If set to true, 'desiredRevision' must be
	// specified.
	UpdateTargetRevision bool `json:"updateTargetRevision,omitempty"`
}

// Describes updates to an Argo CD Application source's Helm parameters.
type ArgoCDHelmParameterUpdates struct {
	Images []ArgoCDHelmImageUpdate `json:"images"`
}

// Describes how to update a Helm parameter to reference a specific version of a container
// image.
type ArgoCDHelmImageUpdate struct {
	// Specifies a key within an Argo CD Application source's Helm parameters that is to be
	// updated.
	Key string `json:"key"`
	// Specifies a new value for the setting within an Argo CD Application source's Helm
	// parameters identified by the 'key' field.
	Value string `json:"value"`
}

// Describes updates to an Argo CD Application source's Kustomize images.
type ArgoCDKustomizeImageUpdates struct {
	Images []ArgoCDKustomizeImageUpdate `json:"images,omitempty"`
}

// Describes how to update a Kustomize image to reference a specific version of a container
// image.
type ArgoCDKustomizeImageUpdate struct {
	// Digest of the image to set. Mutually exclusive with 'tag'.
	Digest string `json:"digest,omitempty"`
	// Specifies a container image name override.
	NewName string `json:"newName,omitempty"`
	// The URL of a container image repository.
	RepoURL string `json:"repoURL"`
	// Tag of the image to set. Mutually exclusive with 'digest'.
	Tag string `json:"tag,omitempty"`
}

type CopyConfig struct {
	// Ignore is a (multiline) string of glob patterns to ignore when copying files. It accepts
	// the same syntax as .gitignore files.
	Ignore string `json:"ignore,omitempty"`
	// InPath is the path to the file or directory to copy.
	InPath string `json:"inPath"`
	// OutPath is the path to the destination file or directory.
	OutPath string `json:"outPath"`
}

type DeleteConfig struct {
	// Path is the path to the file or directory to delete.
	Path string `json:"path"`
	// Strict will cause the directive to fail if the path does not exist.
	Strict bool `json:"strict,omitempty"`
}

type GitClearConfig struct {
	// Path to a working directory of a local repository from which to remove all files,
	// excluding the .git/ directory.
	Path string `json:"path"`
}

type GitCloneConfig struct {
	// The commits, branches, or tags to check out from the repository and the paths where they
	// should be checked out. At least one must be specified.
	Checkout []Checkout `json:"checkout"`
	// Indicates whether to skip TLS verification when cloning the repository. Default is false.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	// The URL of a remote Git repository to clone. Required.
	RepoURL string `json:"repoURL"`
}

type Checkout struct {
	// The branch to checkout. Mutually exclusive with 'commit' and 'tag'. If none of these are
	// specified, the default branch is checked out.
	Branch string `json:"branch,omitempty"`
	// The commit to checkout. Mutually exclusive with 'branch' and 'tag''. If none of these are
	// specified, the default branch is checked out.
	Commit string `json:"commit,omitempty"`
	// Indicates whether a new, empty orphan branch should be created if the branch does not
	// already exist. Default is false.
	Create bool `json:"create,omitempty"`
	// The path where the repository should be checked out.
	Path string `json:"path"`
	// The tag to checkout. Mutually exclusive with 'branch' and 'commit'. If none of these are
	// specified, the default branch is checked out.
	Tag string `json:"tag,omitempty"`
}

type GitCommitConfig struct {
	// The author of the commit.
	Author *Author `json:"author,omitempty"`
	// The commit message. Mutually exclusive with 'messageFromSteps'.
	Message string `json:"message,omitempty"`
	// TODO
	MessageFromSteps []string `json:"messageFromSteps,omitempty"`
	// The path to a working directory of a local repository.
	Path string `json:"path"`
}

// The author of the commit.
type Author struct {
	// The email of the author.
	Email string `json:"email,omitempty"`
	// The name of the author.
	Name string `json:"name,omitempty"`
}

type GitOpenPRConfig struct {
	// Indicates whether a new, empty orphan branch should be created and pushed to the remote
	// if the target branch does not already exist there. Default is false.
	CreateTargetBranch bool `json:"createTargetBranch,omitempty"`
	// Indicates whether to skip TLS verification when cloning the repository. Default is false.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	// Labels to add to the pull request.
	Labels []string `json:"labels,omitempty"`
	// The name of the Git provider to use. Currently only 'github', 'gitlab', 'gitea' and
	// 'azure' are supported. Kargo will try to infer the provider if it is not explicitly
	// specified.
	Provider *Provider `json:"provider,omitempty"`
	// The URL of a remote Git repository to clone.
	RepoURL string `json:"repoURL"`
	// The branch containing the changes to be merged. This branch must already exist and be up
	// to date on the remote.
	SourceBranch string `json:"sourceBranch"`
	// The branch to which the changes should be merged. This branch must already exist and be
	// up to date on the remote.
	TargetBranch string `json:"targetBranch"`
	// The title for the pull request. Kargo generates a title based on the commit messages if
	// it is not explicitly specified.
	Title string `json:"title,omitempty"`
}

type GitPushConfig struct {
	// Indicates whether to push to a new remote branch. A value of 'true' is mutually exclusive
	// with 'targetBranch'. If neither of these is provided, the target branch will be the
	// currently checked out branch.
	GenerateTargetBranch bool `json:"generateTargetBranch,omitempty"`
	// This step implements its own internal retry logic for cases where a push is determined to
	// have failed due to the remote branch having commits that that are not present locally.
	// Each attempt, including the first, rebases prior to pushing. This field configures the
	// maximum number of attempts to push to the remote repository. If not specified, the
	// default is 50.
	MaxAttempts *int64 `json:"maxAttempts,omitempty"`
	// The path to a working directory of a local repository.
	Path string `json:"path"`
	// The target branch to push to. Mutually exclusive with 'generateTargetBranch=true'. If
	// neither of these is provided, the target branch will be the currently checked out branch.
	TargetBranch string `json:"targetBranch,omitempty"`
}

type GitWaitForPRConfig struct {
	// Indicates whether to skip TLS verification when cloning the repository. Default is false.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	// The number of the pull request to wait for.
	PRNumber int64 `json:"prNumber"`
	// The name of the Git provider to use. Currently only 'github', 'gitlab', 'gitea' and
	// 'azure' are supported. Kargo will try to infer the provider if it is not explicitly
	// specified.
	Provider *Provider `json:"provider,omitempty"`
	// The URL of a remote Git repository to clone.
	RepoURL string `json:"repoURL"`
}

type HelmTemplateConfig struct {
	// APIVersions allows a manual set of supported API Versions to be passed when rendering the
	// manifests.
	APIVersions []string `json:"apiVersions,omitempty"`
	// Whether to disable hooks in the rendered manifests.
	DisableHooks bool `json:"disableHooks,omitempty"`
	// Whether to include CRDs in the rendered manifests.
	IncludeCRDs bool `json:"includeCRDs,omitempty"`
	// KubeVersion allows for passing a specific Kubernetes version to use when rendering the
	// manifests.
	KubeVersion string `json:"kubeVersion,omitempty"`
	// Namespace to use for the rendered manifests.
	Namespace string `json:"namespace,omitempty"`
	// OutPath to write the rendered manifests to. If it points to a .yaml or .yml file, the
	// rendered manifests will be written to that file. If it points to a directory, the
	// rendered manifests will be written to this directory joined with the chart name.
	OutPath string `json:"outPath"`
	// Path at which the Helm chart can be found.
	Path string `json:"path"`
	// ReleaseName to use for the rendered manifests.
	ReleaseName string `json:"releaseName"`
	// Whether to skip tests when rendering the manifests.
	SkipTests bool `json:"skipTests,omitempty"`
	// Whether to use the release name in the output path (instead of the chart name). This only
	// has an effect if outPath is set to a directory.
	UseReleaseName bool `json:"useReleaseName,omitempty"`
	// ValuesFiles to use for rendering the Helm chart.
	ValuesFiles []string `json:"valuesFiles,omitempty"`
}

type HelmUpdateChartConfig struct {
	// A list of chart dependencies which should receive updates.
	Charts []Chart `json:"charts"`
	// The path at which the umbrella chart with the dependency can be found.
	Path string `json:"path"`
}

type Chart struct {
	// The name of the subchart, as defined in `Chart.yaml`.
	Name string `json:"name"`
	// The repository of the subchart, as defined in `Chart.yaml`. It also supports OCI charts
	// using `oci://`.
	Repository string `json:"repository"`
	// The version of the subchart to update to.
	Version string `json:"version"`
}

type HTTPConfig struct {
	// The body of the HTTP request.
	Body string `json:"body,omitempty"`
	// An expression to evaluate to determine if the request failed.
	FailureExpression string `json:"failureExpression,omitempty"`
	// Headers to include in the HTTP request.
	Headers []HTTPHeader `json:"headers,omitempty"`
	// Whether to skip TLS verification when making the request. (Not recommended.)
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	// The HTTP method to use for the request.
	Method string `json:"method,omitempty"`
	// Outputs to extract from the HTTP response.
	Outputs []HTTPOutput `json:"outputs,omitempty"`
	// Query parameters to include in the HTTP request.
	QueryParams []HTTPQueryParam `json:"queryParams,omitempty"`
	// An expression to evaluate to determine if the request was successful.
	SuccessExpression string `json:"successExpression,omitempty"`
	// The maximum time to wait for the request to complete. If not specified, the default is 10
	// seconds.
	Timeout string `json:"timeout,omitempty"`
	// The URL to send the HTTP request to.
	URL string `json:"url"`
}

type HTTPHeader struct {
	// The name of the header.
	Name string `json:"name"`
	// The value of the header.
	Value string `json:"value"`
}

type HTTPOutput struct {
	// An expression to evaluate to extract the output from the HTTP response.
	FromExpression string `json:"fromExpression"`
	// The name of the output.
	Name string `json:"name"`
}

type HTTPQueryParam struct {
	// The name of the query parameter.
	Name string `json:"name"`
	// The value of the query parameter.
	Value string `json:"value"`
}

type JSONUpdateConfig struct {
	// The path to a JSON file.
	Path string `json:"path"`
	// A list of updates to apply to the JSON file.
	Updates []JSONUpdate `json:"updates"`
}

type JSONUpdate struct {
	// The key whose value needs to be updated. For nested values, use a JSON dot notation path.
	Key string `json:"key"`
	// The new value for the specified key.
	Value interface{} `json:"value"`
}

type KustomizeBuildConfig struct {
	// OutPath is the file path to write the built manifests to.
	OutPath string `json:"outPath"`
	// Path to the directory containing the Kustomization file.
	Path string `json:"path"`
	// Plugin contains configuration for customizing the behavior of builtin Kustomize plugins.
	Plugin *Plugin `json:"plugin,omitempty"`
}

// Plugin contains configuration for customizing the behavior of builtin Kustomize plugins.
type Plugin struct {
	// Helm contains configuration for inflating a Helm chart.
	Helm *Helm `json:"helm,omitempty"`
}

// Helm contains configuration for inflating a Helm chart.
type Helm struct {
	// APIVersions allows a manual set of supported API versions to be passed when inflating a
	// Helm chart.
	APIVersions []string `json:"apiVersions,omitempty"`
	// KubeVersion allows for passing a specific Kubernetes version to use when inflating a Helm
	// chart.
	KubeVersion string `json:"kubeVersion,omitempty"`
}

type KustomizeSetImageConfig struct {
	// Images is a list of container images to set or update in the Kustomization file. When
	// left unspecified, all images from the Freight collection will be set in the Kustomization
	// file. Unless there is an ambiguous image name (for example, due to two Warehouses
	// subscribing to the same repository), which requires manual configuration.
	Images []Image `json:"images"`
	// Path to the directory containing the Kustomization file.
	Path string `json:"path"`
}

type Image struct {
	// Digest of the image to set in the Kustomization file. Mutually exclusive with 'tag'.
	Digest string `json:"digest,omitempty"`
	// Image name of the repository from which to pick the version. This is the image name Kargo
	// is subscribed to, and produces Freight for.
	Image string `json:"image"`
	// Name of the image (as defined in the Kustomization file).
	Name string `json:"name,omitempty"`
	// NewName for the image. This can be used to rename the container image name in the
	// manifests.
	NewName string `json:"newName,omitempty"`
	// Tag of the image to set in the Kustomization file. Mutually exclusive with 'digest'.
	Tag string `json:"tag,omitempty"`
}

type YAMLUpdateConfig struct {
	// The path to a YAML file.
	Path string `json:"path"`
	// A list of updates to apply to the YAML file.
	Updates []YAMLUpdate `json:"updates"`
}

type YAMLUpdate struct {
	// The key whose value needs to be updated. For nested values, use a YAML dot notation path.
	Key string `json:"key"`
	// The new value for the specified key.
	Value string `json:"value"`
}

// The name of the Git provider to use. Currently only 'github', 'gitlab', 'gitea' and
// 'azure' are supported. Kargo will try to infer the provider if it is not explicitly
// specified.
type Provider string

const (
	Azure  Provider = "azure"
	Gitea  Provider = "gitea"
	Github Provider = "github"
	Gitlab Provider = "gitlab"
)
