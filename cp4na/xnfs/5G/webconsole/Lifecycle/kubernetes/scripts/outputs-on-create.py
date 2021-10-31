def getOutputs(keg, props, resultBuilder, log, *args, **kwargs):
	found, webconsole = keg.objects.get('5g.ibm.com/v1alpha1', 'Webconsole', props['system_properties']['resource_subdomain'], namespace=props['deployment_location']['properties']['defaultObjectNamespace'])
	if not found:
		return resultBuilder.failed('Could not find Webconsole')

	status = webconsole['status']
	if 'outputs' in status:
		outputs = status['outputs']
		# resultBuilder.setOutput('x1', str(outputs))
		# resultBuilder.setOutput('ip_address', outputs["ip_address"])
		# resultBuilder.setOutput('ip_address', outputs['ip_address'])
		for key in outputs:
			resultBuilder.setOutput(key, outputs[key])
