app = angular.module "Hello", ["ngResource"]

app.factory "Users", ($resource, $q) ->
	$resource "/data/:id"

app.controller "HelloCtrl", (Users) ->
	controller =
		body: "body..."
		title: "HelloCtrl title"
		pickUser: (id) ->
			resp = Users.get id: id
			resp.$promise.then (user) ->
				controller.user = user

	fetchUsers = ->
		Users.query!.$promise.then (users) ->
			controller.users = users

	controller.delete = ->
		delete controller.user
		Users.delete!.$promise.then fetchUsers

	controller.create = ->
		Users.save!.$promise.then fetchUsers

	fetchUsers!
	controller

