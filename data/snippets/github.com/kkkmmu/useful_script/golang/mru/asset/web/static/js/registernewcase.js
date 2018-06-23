var next = 1;
$(".add-more").click(function(e){
	e.preventDefault();
	var addto = "#field" + next;
	next = next + 1;
	var newIn = '<br /><br /><input autocomplete="off" class="span3" id="field' + next + '" name="field' + next + '" type="text" data-provide="typeahead" data-items="8">';
	var newInput = $(newIn);
	$(addto).after(newInput);
	$("#field" + next).attr('data-source',$(addto).attr('data-source'));
	$("#count").val(next);  
});

$("#remove").click(function(){
	$("#table").remove();
});

$("#add").click(function(){
	$("#table").clone().appendTo("#div1");
});

var a=1;
$("#add_assertion").click(function(){
	$('#assert'+a).html("<td>"+ (a+1) +"</td><td><input name='mode"+a+"' type='text' placeholder='enable' class='form-control input-md'  /> </td><td><input  name='cli"+a+"' type='text' placeholder='show running-config'  class='form-control input-md'></td> <td><input  name='expect"+a+"' type='text' placeholder='br100[[:space:]]+up'  class='form-control input-md'></td>");

	$('#assertion_tab').append('<tr id="assert'+(a+1)+'"></tr>');
	a++; 
});
$("#delete_assertion").click(function(){
	if(a>1){
		$("#assert"+(a-1)).html('');
		a--;
	}
});


var c=1;
$("#add_command").click(function(){
	$('#command'+c).html("<td>"+ (c+1) +"</td><td><input name='mode"+c+"' type='text' placeholder='enable' class='form-control input-md'  /> </td><td><input  name='cli"+c+"' type='text' placeholder='show running-config'  class='form-control input-md'></td>");

	$('#command_tab').append('<tr id="command'+(c+1)+'"></tr>');
	c++; 
});
$("#delete_command").click(function(){
	if(c>1){
		$("#command"+(c-1)).html('');
		c--;
	}
});

var task=1;

$("#deltask").click(function(){
		$("#deltask").remove();
});

$("#addtask").click(function(){
	$("#task0").clone().appendTo("#tasks");
});

var routine=1;

$("#delroutine").click(function(){
		$("#delroutine").remove();
});

$("#addroutine").click(function(){
	$("#routine0").clone().appendTo("#routine0");
});
