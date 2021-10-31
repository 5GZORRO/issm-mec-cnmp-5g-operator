def checkReady(keg, props, resultBuilder, log, *args, **kwargs):
    namespace = props['deployment_location']['properties']['defaultObjectNamespace']
    name = props['system_properties']['resource_subdomain']
    found, mongo = keg.objects.get('5g.ibm.com/v1alpha1', 'Mongo', name, namespace=namespace)
    if not found:
        return resultBuilder.failed(f'Could not find Mongo {name} in namespace {namespace}')

    metadata = mongo['metadata']
    status = mongo.get('status', None)
    if status is not None:
        ready = False
        if 'generation' in metadata:
            generation = metadata['generation']
            if 'observedGeneration' in status:
                observedGeneration = status['observedGeneration']
                if observedGeneration >= generation:
                    ready = True

        if 'conditions' in status:
            ready_conditions = [condition for condition in status['conditions'] if condition['type'] == 'Reconciled' and condition['status'] == 'True']
            if ready and len(ready_conditions) > 0:
                return resultBuilder.ready()
            else:
                return resultBuilder.notReady()
        else:
            return resultBuilder.notReady()
    else:
        return resultBuilder.notReady()
