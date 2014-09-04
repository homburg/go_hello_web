app = window.angular.module "Hello", ["ngResource"]

ng = {}

app.factory "Users", ($resource) ->
  /**
   * @type {ng.Resource}
   */
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

  /**
   * @param {number} text
   */
  setText = (text) ->
    controller.body = text

  setText 42

  fetchUsers!
  controller

