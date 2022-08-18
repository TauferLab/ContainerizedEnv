// var devel=1;
var devel=0;
var widgets = require('jupyter-js-widgets');
var _ = require('underscore');

// Custom Model. Custom widgets models must at least provide default values
// for model attributes, including
//
//  - `_view_name`
//  - `_view_module`
//  - `_view_module_version`
//
//  - `_model_name`
//  - `_model_module`
//  - `_model_module_version`
//
//  when different from the base class.

// When serialiazing the entire widget state for embedding, only values that
// differ from the defaults will be specified.
var HelloModel = widgets.DOMWidgetModel.extend({
    defaults: _.extend(_.result(this, 'widgets.DOMWidgetModel.prototype.defaults'), {
	_model_name : 'HelloModel',
	_view_name : 'HelloView',
	_model_module : 'ipywe',
	_view_module : 'ipywe',
	_model_module_version : '0.1.0',
	_view_module_version : '0.1.0',
	value : 'Hello World'
    })
});


// Custom View. Renders the widget model.
var HelloView = widgets.DOMWidgetView.extend({
    render: function() {
	this.value_changed();
	this.model.on('change:value', this.value_changed, this);
    },

    value_changed: function() {
	this.el.textContent = this.model.get('value');
    }
});


if (devel) {
    require.undef("example");
    define("example", ["jupyter-js-widgets"], function(widgets) {
	return {
	    HelloView: HelloView,
	    HelloModel: HelloModel
	}
    });
} else {
    module.exports = {
	HelloModel : HelloModel,
	HelloView : HelloView
    };
}


