app = angular.module "Hello", ["ngResource"]

app.factory "Users", ($resource) ->
	res = $resource "/data/:id"

	/**
	 * @param {number} id
	 */
	get: (id) -> res.get id: id .$promise
	query: -> res.query!.$promise
	delete: -> res.delete!.$promise
	create: -> res.save!.$promise


app.controller "HelloCtrl", (Users) ->
	controller =
		body: "body..."
		title: "HelloCtrl title"
		pickUser: (id) ->
			resp = Users.get id
			resp.then (user) ->
				controller.user = user

	fetchUsers = ->
		Users.query!.then (users) ->
			controller.users = users

	controller.delete = ->
		delete controller.user
		Users.delete!.then fetchUsers

	controller.create = ->
		Users.create!.then fetchUsers

	fetchUsers!
	controller

