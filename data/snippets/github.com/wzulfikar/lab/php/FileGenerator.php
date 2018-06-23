<?php

class FileGenerator{
	public $vars;
	public $template;
	public function __construct($template, array $vars = [])
	{
		if ( ! is_string($template) ) 
			throw new Exception('Type of $template should be string instead of ' . gettype($template));
		$this->template = $template;
		$this->setVars($vars);
	}
	public function setVars(array $vars){
		$this->vars = $vars;
	}
	public function parse(){
		return $this->replace_template_vars($this->template, $this->vars);
	}
	public function put($output){
		file_put_contents($output, $this->parse());
	}
	private function replace_template_vars($template, array $vars){
	    foreach ($vars as $var_name_to_replace => $new_var_name) {
	        $template = str_replace('{{' . $var_name_to_replace . '}}', $new_var_name, $template);
	    }
	    return $template;
	}
}
$template = file_get_contents('ModelTemplate.php');
$vars = [
	'namespace' => 'App\Polymorphic\Likeable',
	'modelName' => 'Like',
	'polymorphicName' => 'Likeable'
];
$template = new FileGenerator($template,$vars);
// check what the generated file will look like
var_dump($template->parse());
// store the file
$template->put($vars['modelName'] . '.php');
// sample of ModelTemplate.php (input)
// <?php
// 
// namespace {{namespace}};
// 
// use Illuminate\Database\Eloquent\Model;
// 
// class {{modelName}} extends Model
// {
//   /**
//    * Get all of the owning {{polymorphicName}} models.
//    */
//   public function {{polymorphicName}}()
//   {
//       return $this->morphTo();
//   }
// }
// 
// The output will be `Like.php`
//
// <?php
// 
// namespace App\Polymorphic\Likeable;
// 
// use Illuminate\Database\Eloquent\Model;
// 
// class Like extends Model
// {
//   /**
//    * Get all of the owning Likeable models.
//    */
//   public function Likeable()
//   {
//       return $this->morphTo();
//   }
// }
