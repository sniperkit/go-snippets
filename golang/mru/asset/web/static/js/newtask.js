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
	alert("Hello world")
	//alert($("#add_command").prev().name())
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

var routine=1;

$("#delroutine").click(function(){
		$("#delroutine").remove();
});

$("#addroutine").click(function(){
	$("#routine0").clone().appendTo("#routine0");
});

$(#"continue").onclick(function() {
	alert("Continue")
}

$(#"finish").click(function() {
	alert("finish")
}

$().ready(function() {
	$("#newtaskform").validate();
});
