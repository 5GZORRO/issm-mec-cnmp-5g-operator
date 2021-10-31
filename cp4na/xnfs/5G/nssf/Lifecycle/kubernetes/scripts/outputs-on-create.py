def getOutputs(keg, props, resultBuilder, log, *args, **kwargs):
	found, nssf = keg.objects.get('5g.ibm.com/v1alpha1', 'Nssf', props['system_properties']['resource_subdomain'], namespace=props['deployment_location']['properties']['defaultObjectNamespace'])
	if not found:
		return resultBuilder.failed('Could not find Nssf')

	status = nssf['status']
	if 'outputs' in status:
		outputs = status['outputs']
		for key in outputs:
			resultBuilder.setOutput(key, outputs[key])
