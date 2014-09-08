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

app.factory "Connection", ->
  w = null
  connect = ->
    if !w?
      w = new WebSocket("ws://localhost:3000/socket")
    w
  {
    connect: (callback) ->
      w = connect!
      w.onmessage = callback

    connectJson: (callback) ->
      w = connect!
      w.onmessage = (event) ->
        callback JSON.parse event.data
  }

app.controller "HelloCtrl", ($scope, Users, Connection) ->
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
    Users.delete!

  controller.create = ->
    Users.create!

  Connection.connectJson (data) ->
    $scope.$apply ->
      controller.users = data

  /**
   * @param {number} text
   */
  setText = (text) ->
    controller.body = text

  setText 42

  fetchUsers!
  controller
