def getOutputs(keg, props, resultBuilder, log, *args, **kwargs):
  namespace = props['deployment_location']['properties']['defaultObjectNamespace']
  name = props['system_properties']['resource_subdomain']
  found, upf = keg.objects.get('5g.ibm.com/v1alpha1', 'Upf', name, namespace=namespace)
  if not found:
    return resultBuilder.failed(f'Could not find Upf {name} in {namespace}')

  status = upf['status']
  if 'outputs' in status:
    outputs = status['outputs']
    for key in outputs:
      resultBuilder.setOutput(key, outputs[key])
