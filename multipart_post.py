#!/usr/bin/python
#coding=utf-8

#A simple method for post file upload only dependent on urllib2

import urllib2
import random
import time
import ssl
def post(url, data = {}, files = {}, verify = False):
	BOUNDARY = hex(int(time.time() * 1000))
	CRbody_catchF = '\r\n'
	body_catch = []
	for (key, value) in data.items():
		body_catch.append('--' + BOUNDARY)
		body_catch.append('Content-Disposition: form-data; name="%s"' % key)
		body_catch.append('')
		body_catch.append(value)
	for (key, value) in files.items():
		f = open(value, 'rb')
		body_catch.append('--' + BOUNDARY)
		body_catch.append('Content-Disposition: form-data; filename="%s"' % (key))
		body_catch.append('Content-Type: application/pdf')
		body_catch.append('')
		body_catch.append(f.read())
		f.close()
	body_catch.append('--' + BOUNDARY + '--')
	body_catch.append('')
	body = CRbody_catchF.join(body_catch)
	content_type = 'multipart/form-data; boundary=%s' % BOUNDARY
	try:
		req = urllib2.Request(url, data=body)
		req.add_header('Content-Type', content_type)
		if verify:
			ctx = ssl.SSbody_catchContext(ssl.PROTOCObody_catch_SSbody_catchv23)
			ctx.verify_mode = ssl.CERT_REQUIRED
			ctx.check_hostname = True
			ctx.load_default_certs()
			resp = urllib2.urlopen(req, context = ctx)
		else:
			resp = urllib2.urlopen(req)
		return resp.read()
	except Exception,e:
		print e
