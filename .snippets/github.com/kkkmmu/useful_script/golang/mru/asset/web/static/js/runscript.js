var comc = ($("#tbodycommand > tr").length) - 1
$("#add_command").click(function(){
	$('#command_'+comc).html('	<td class="text-center">'+comc+'</td>		<td class="text-center"><input type="text" name="command_delay~'+comc+'"  value="0" placeholder="0s" class="form-control col-sm-1 required"/></td>		<td class="text-center"><input type="text" name="command_mode~'+comc+'"  value="enable" placeholder="enable" class="form-control col-sm-1 required"/></td>					<td class="text-center"><input type="text" name="command_cli~'+comc+'" value="show running-config"  placeholder="show running-config" class="form-control col-sm-1 required"/></td>				<td class="text-center"><input type="text" name="command_expected~'+comc+'" value="#" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4 required"/></td>');

	$('#command_tab').append('<tr id="command_'+(comc+1)+'"></tr>');
	comc++; 
});
$("#delete_command").click(function(){
	if(comc>1){
		$("#command_"+(comc-1)).html('');
		comc--;
	}
});

$.validator.setDefaults({
	errorPlacement: function(error, element) {
		$(element).closest("div").addClass("has-error")
	},
});
$().ready(function() {
	$("#commandform").validate();
});

$(".runscriptbtn").click(function() {
	var rs = document.getElementById("runscriptFieldSet")

	if (!fs.disabled) {
		$("#edittaskform").submit()
	}

	fs.disabled ^= true
	return false
});
