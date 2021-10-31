def run(keg, props, resultBuilder, log, *args, **kwargs):
    namespace = props['deployment_location']['properties']['defaultObjectNamespace']
    name = props['system_properties']['resource_subdomain']
    found, amf = keg.objects.get('5g.ibm.com/v1alpha1', 'Amf', name, namespace=namespace)
    if not found:
        return resultBuilder.failed(f'Could not find Amf {name} in namespace {namespace}')

    spec = amf['spec']
    spec['targetState'] = started
    metadata = amf['metadata']
    status = amf.get('status', None)
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

