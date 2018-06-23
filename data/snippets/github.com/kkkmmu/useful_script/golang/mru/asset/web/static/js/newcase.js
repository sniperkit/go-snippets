var d=1;
$("#add_dut").click(function(){
	$('#row'+d).html("<td>"+ (d+1) +"</td><td><input  name='device~"+d+"' type='text' placeholder='V8500_SFU'  class='form-control required input-md'></td><td><input name='username~"+d+"' type='text' placeholder='admin' class='form-control required input-md'  /> </td><td><input name='password~"+d+"' type='text' placeholder='Dasan123456' class='form-control input-md'  /> ");

	$('#dut_tab').append('<tr id="row'+(d+1)+'"></tr>');
	d++; 
});
$("#del_dut").click(function(){
	if(d>1){
		$("#row"+(d-1)).html('');
		d--;
	}
});

$.validator.setDefaults({
	errorPlacement: function(error, element) {
		$(element).closest("div").addClass("has-error")
	},
});

$().ready(function() {
	$("#caseform").validate();
});
