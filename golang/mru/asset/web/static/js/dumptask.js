var prec = ($("#tbodyprecondition > tr").length) - 1
$("#precondition_add_assert").click(function(){
	$('#preconditionassert_'+prec).html('	<td class="text-center">'+prec+'</td>											<td class="text-center"><input type="text" name="preconditionassertdut~'+prec+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="preconditionassertmode~'+prec+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="preconditionassertcli~'+prec+'" placeholder="show runngin-config" class="form-control col-sm-4 required"/></td>											<td class="text-center"><input type="text" name="preconditionassertexpected~'+prec+'" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4 required"/></td>');

	$('#preconditionassert_tab').append('<tr id="preconditionassert_'+(prec+1)+'"></tr>');
	prec++; 
});
$("#precondition_delete_assert").click(function(){
	if(prec>1){
		$("#preconditionassert_"+(prec-1)).html('');
		prec--;
	}
});

var	posc = ($("#tbodypostcondition > tr").length) - 1
$("#postcondition_add_assert").click(function(){
	$('#postconditionassert_'+posc).html('	<td class="text-center">'+posc+'</td>											<td class="text-center"><input type="text" name="postconditionassertdut~'+posc+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="postconditionassertmode~'+posc+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="postconditionassertcli~'+posc+'" placeholder="show runngin-config" class="form-control col-sm-4 required"/></td>											<td class="text-center"><input type="text" name="postconditionassertexpected~'+posc+'" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4 required"/></td>');

	$('#postconditionassert_tab').append('<tr id="postconditionassert_'+(posc+1)+'"></tr>');
});
$("#postcondition_delete_assert").click(function(){
	alert(posc)
	if(posc>1){
		$("#postconditionassert_"+(posc-1)).html("");
		posc--;
	}
});


var	comc = ($("#tbodyroutine > tr").length) - 1
$("#routine_add_command").click(function(){
	$('#routine_command_'+comc).html('	<td class="text-center">'+comc+'</td>											<td class="text-center"><input type="text" name="routine_command_dut~'+comc+'"  placeholder="DUT1" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="routine_command_mode~'+comc+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="routine_command_cli~'+comc+'" placeholder="show runngin-config" class="form-control col-sm-4 required"/></td>											<td class="text-center"><input type="text" name="routine_command_expected~'+comc+'" value="#" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4 required"/></td>');

	$('#routine_command_tab').append('<tr id="routine_command_'+(comc+1)+'"></tr>');
	comc++; 
});
$("#routine_delete_command").click(function(){
	if(comc>1){
		$("#routine_command_"+(comc-1)).html('');
		comc--;
	}
});

var	clcomc = ($("#tbodyclear > tr").length) - 1
$("#clear_add_command").click(function(){
	$('#clear_command_'+clcomc).html('	<td class="text-center">'+clcomc+'</td>											<td class="text-center"><input type="text" name="clear_command_dut~'+clcomc+'"  placeholder="DUT1" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="clear_command_mode~'+clcomc+'"  placeholder="enable" class="form-control col-sm-1 required"/></td>											<td class="text-center"><input type="text" name="clear_command_cli~'+clcomc+'" placeholder="show runngin-config" class="form-control col-sm-4 required"/></td>											<td class="text-center"><input type="text" name="clear_command_expected~'+clcomc+'" value="#" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4 required"/></td>');

	$('#clear_command_tab').append('<tr id="clear_command_'+(clcomc+1)+'"></tr>');
	clcomc++; 
});

$("#clear_delete_command").click(function(){
	if(clcomc>1){
		$("#clear_command_"+(clcomc-1)).html('');
		clcomc--;
	}
});

$.validator.setDefaults({
	errorPlacement: function(error, element) {
		$(element).closest("div").addClass("has-error")
	},
});

$().ready(function() {
	$("#edittaskform").validate();
});

$(".editbtn").click(function() {
	var fs = document.getElementById("taskFieldSet")

	if (!fs.disabled) {
		$("#edittaskform").submit()
	}

	fs.disabled ^= true
	return false
});
