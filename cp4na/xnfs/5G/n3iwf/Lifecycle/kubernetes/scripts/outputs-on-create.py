def getOutputs(keg, props, resultBuilder, log, *args, **kwargs):
	found, n3iwf = keg.objects.get('5g.ibm.com/v1alpha1', 'N3iwf', props['system_properties']['resource_subdomain'], namespace=props['deployment_location']['properties']['defaultObjectNamespace'])
	if not found:
		return resultBuilder.failed('Could not find N3iwf')

	status = n3iwf['status']
	if 'outputs' in status:
		outputs = status['outputs']
		for key in outputs:
			resultBuilder.setOutput(key, outputs[key])
