/**
 * @type {Object}
 * @const
 */
var ng = {};

/**
 * @typedef {{
 *   query: ng.ResourceFun,
 *   get: ng.ResourceGetFun,
 *   delete: ng.ResourceFun,
 *   save: ng.ResourceFun,
 *   $promise: Promise
 *  }}
 */
ng.Resource;

/**
 * @typedef {function(): ng.Resource}
 */
ng.ResourceFun;

/**
 * @typedef {function(Object): ng.Resource}
 */
ng.ResourceGetFun;

/**
 * @typedef {{then: function()}}
 */
ng.Promise;
