document.getElementsByClassName = Element.prototype.getElementsByClassName = function(class_names) {
	class_names = (' ' + class_names)
	.replace(/[~!@$%^&*()_+\-=,./';:"?><[\]{}|`#]/g, '\\$&')
	.replace(/\s*(\s|$)/g, '$1')
	.replace(/\s/g, '.');
	return this.querySelectorAll(class_names);
};