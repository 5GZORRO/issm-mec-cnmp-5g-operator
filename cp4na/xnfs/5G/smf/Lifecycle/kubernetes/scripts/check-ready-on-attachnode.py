def checkReady(keg, props, resultBuilder, log, *args, **kwargs):
    found, attach_node = keg.objects.get('5g.ibm.com/v1alpha1', 'Transition', props['system_properties']['resource_subdomain'] + '-' + props['nodeId'] + "-attachnode", namespace=props['deployment_location']['properties']['defaultObjectNamespace'])
    if not found:
        return resultBuilder.failed('Could not find attach_node transition')

    metadata = attach_node['metadata']
    status = attach_node.get('status', None)

    if status is not None:
        ready = False
        if 'generation' in metadata:
            generation = metadata['generation']
            if 'observedGeneration' in status:
                observedGeneration = status['observedGeneration']
                if observedGeneration >= generation:
                    ready = True

        if 'conditions' in status:
            # TODO check observedGeneration
            ready_conditions = [condition for condition in status['conditions'] if condition['type'] == 'Ready' and condition['status'] == 'True']
            if ready and len(ready_conditions) > 0:
                return resultBuilder.ready()
            else:
                return resultBuilder.notReady()
        else:
            return resultBuilder.notReady()
    else:
        return resultBuilder.notReady()
