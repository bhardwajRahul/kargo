import { create } from '@bufbuild/protobuf';
import { createConnectQueryKey, createProtobufSafeUpdater } from '@connectrpc/connect-query';

import { queryClient } from '@ui/config/query-client';
import { listAnalysisTemplates } from '@ui/gen/api/service/v1alpha1/service-KargoService_connectquery';
import { ListAnalysisTemplatesResponseSchema } from '@ui/gen/api/service/v1alpha1/service_pb';
import { AnalysisTemplate } from '@ui/gen/api/stubs/rollouts/v1alpha1/generated_pb';

export default {
  add: (project: string, templates: AnalysisTemplate[]) => {
    queryClient.setQueriesData(
      {
        queryKey: createConnectQueryKey({
          schema: listAnalysisTemplates,
          cardinality: 'finite',
          input: { project }
        }),
        exact: false
      },
      createProtobufSafeUpdater(listAnalysisTemplates, (prev) => {
        let newTemplates = [...(prev?.analysisTemplates || [])];

        if (templates?.length > 0) {
          newTemplates = newTemplates.concat(templates);
        }

        return create(ListAnalysisTemplatesResponseSchema, {
          analysisTemplates: newTemplates
        });
      })
    );
  }
};
