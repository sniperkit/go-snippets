define(["jquery", "swig", "bootstrap"], function($, swig){
	function processTemplateAsString(id){
		return swig.render(document.getElementById(id).innerHTML, { });
	}
	$("body").on("blur click dblclick focus", function(e){
		console.log(e.type, e.target);
	});
	$(".right-pane .btn").popover({
		content: processTemplateAsString("user-action-popover"),
		html: true,
		trigger: "focus",
		template: '<div class="popover" role="tooltip"><div class="arrow"></div><div class="popover-content"></div></div>'
	});
});
