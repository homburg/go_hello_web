app = angular.module "Hello", ["ngResource"]

app.factory "Users", ($resource, $q) ->
	res = $resource "/data/:id"
	cache = {}

	query: ->
		res.query!
	get: (id) ->
		if id of cache
			cache[id]
		else
			resResp = res.get {id: id}
			cache[id] = resResp
			resResp

app.controller "HelloCtrl", (Users) ->
	controller =
		body: "body..."
		title: "HelloCtrl title"
		pickUser: (id) ->
			resp = Users.get id
			resp.$promise.then (user) ->
				controller.user = user

	Users.query!.$promise.then (users) ->
		controller.users = users

	controller

