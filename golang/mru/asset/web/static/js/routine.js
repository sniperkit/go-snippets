$("#add_command").click(function() {
	alert("Hello world");
});

$("#del_command").click(function() {
	alert("Goodbye world");
});

var precondition = 1;
$("#add_precondition").click(function(){
	alert("Add PreCondition "+precondition);
	$('#row'+precondition).html("<td>"+ (precondition+1) +"</td><td><input name='dut"+precondition+"' type='text' placeholder='DUT1' class='form-control input-md'  /> </td><td><input  name='device"+precondition+"' type='text' placeholder='V8500_SFU'  class='form-control input-md'></td>");

	$('#dut_tab').append('<tr id="row'+(precondition+1)+'"></tr>');
	precondition++; 
});
$("#del_precondition").click(function(){
	alert("Del PreCondition "+precondition);
	if(precondition>1){
		$("#row"+(precondition-1)).html('');
		precondition--;
	}
});

var routine = 1;
$("#add_routine").click(function(){
	alert("Add Routine "+routine);
	$("#routine"+routine).html('<fieldset><div class="row-fluid col-xs-offset-2"><legend><<<<<'+routine+'>>>>></legend><div class="form-group">	<label class="col-sm-2 control-label" for="textinput">Name</label><div class="col-sm-4"><input class="form-control" id="{{.}}{{`_'+routine'+_`}}{{`name`}}" name="{{.}}{{`_'+routine+'_`}}{{`name`}}" type="text" placeholder="Create VLAN Interface"/></div></div><div class="form-group"><label class="col-sm-2 control-label" for="textinput">DUT</label><div class="col-sm-4"><input class="form-control" id="{{.}}{{`_'+routine+'_`}}{{`dut`}}" name="{{.}}{{`_'+routine+'_`}}{{`dut`}}" type="text" placeholder="DUT1"/></div></div><div class="form-group"><label class="col-sm-2 control-label" for="routine_description" >Command</label><div class="container"><div class="row"><div class="col-sm-6 column"><table class="table table-bordered table-hover" id="{{.}}{{`_'+routine+'_`}}{{`command_tab`}}"><thead><tr ><th class="text-center col-sm-1">#</th><th class="text-center col-sm-1">Mode</th><th class="text-center col-sm-4">CLI</th>	</tr></thead><tbody><tr id="{{.}}{{`_'+routine+'_`}}{{`command`}}{{`_0_`}}"><td class="text-center">0</td><td class="text-center"><input type="text" name="{{.}}{{`_'+routine+'_`}}{{`command`}}{{`_0_`}}{{`mode`}}"  placeholder="enable" class="form-control col-sm-1"/></td><td class="text-center"><input type="text" name="{{.}}{{`_'+routine+'_`}}{{`command`}}{{`_0_`}}{{`cli`}}" placeholder="show runngin-config" class="form-control col-sm-4"/></td></tr><tr id="{{.}}{{`_'+0+'`}}{{`command`}}{{`_1_`}}"></tr></tbody></table><a id="add_command" class="pull-left">Add Commnad</a><a id="delete_command" class="pull-right">Delete Command</a></div></div></div></div><!-- from group--><div class="form-group"><label class="col-sm-2 control-label" for="routine_description">Assertion</label><div class="container"><div class="row"><div class="col-sm-10 column"><table class="table table-bordered table-hover" id="{{.}}{{`_'+routine+'_`}}{{`assertion_tab`}}"><thead><tr ><th class="text-center col-sm-1">#</th><th class="text-center col-sm-1">Mode</th><th class="text-center col-sm-4">CLI</th><th class="text-center col-sm-4">Expected</th></tr></thead><tbody><tr id="{{.}}{{`_'+routine+'_`}}{{`assert`}}{{`_0_`}}"><td class="text-center">1</td><td class="text-center"><input type="text" name="{{.}}{{`_'+routine+'_`}}{{`assert`}}{{`_0_`}}{{`mode`}}"  placeholder="enable" class="form-control col-sm-1"/></td><td class="text-center"><input type="text" name="{{.}}{{`_'+routine+'_`}}{{`assert`}}{{`_0_`}}{{`cli`}}" placeholder="show runngin-config" class="form-control col-sm-4"/></td><td class="text-center"><input type="text" name="{{.}}{{`_'+routine'+_`}}{{`assert`}}{{`_0_`}}{{`expected`}}" placeholder="br1000[[_space_]]+up" class="form-control col-sm-4"/></td></tr><tr id="{{.}}{{`_'+routine+'_`}}{{`assert`}}{{`_1_`}}"></tr></tbody></table><a id="add_assertion" class="pull-left">Add Assertion</a><a id="delete_assertion" class="pull-right">Delete Assertion</a></div></div></div></div><div class="form-group"><label for="routine_description" class="col-sm-2 control-label">Description</label><div class="container"><div class="row"><div class="col-sm-6 column"><textarea id="{{.}}{{`_'+routine+'_`}}{{`routine_description`}}" name="{{.}}{{`_'+routine+'_`}}{{`routine_description`}}" class="form-control col-sm-6" rows="5"></textarea></div></div></div></div><div class="form-group"><div class="col-sm-10 column"></div></div></div></fieldset>');
	$("#routine_tab").append('<div class="routine" id="routine"+routine></div>');
);
	routine++;
})

$("#del_routine").click(function(){
	alert("Del Routine "+routine);
	routine--;
})
var postcondition = 1;
$("#add_postcondition").click(function(){
	alert("Add PostCondtion "+postcondition);
	postcondition++;
})

$("#del_postcondition").click(function(){
	alert("Del PostCondtion "+postcondition);
	postcondition--;
})
