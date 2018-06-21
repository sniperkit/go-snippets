Jquery在异步提交方面封装的很好，直接用AJAX非常麻烦，Jquery大大简化了我们的操作，不用考虑浏览器的诧异了。

推荐一篇不错的jQuery Ajax 实例文章，忘记了可以去看看，

地址为：http://www.cnblogs.com/yeer/archive/2009/07/23/1529460.html 和 http://www.w3school.com.cn/jquery/

$.post、$.get是一些简单的方法，如果要处理复杂的逻辑，还是需要用到jQuery.ajax()

一、$.ajax的一般格式

	$.ajax({
		type: 'POST',
		url: url ,
		data: data ,
		success: success ,
		dataType: dataType
	});

二、$.ajax的参数描述

参数 描述

url	必需。规定把请求发送到哪个 URL。
data	可选。映射或字符串值。规定连同请求发送到服务器的数据。
success(data, textStatus, jqXHR)	可选。请求成功时执行的回调函数。
dataType	
可选。规定预期的服务器响应的数据类型。

默认执行智能判断（xml、json、script 或 html）。

三、$.ajax需要注意的一些地方：

1.data主要方式有三种，html拼接的，json数组，form表单经serialize()序列化的；通过dataType指定，不指定智能判断。

2.$.ajax只提交form以文本方式，如果异步提交包含<file>上传是传过不过去,需要使用jquery.form.js的$.ajaxSubmit

四、$.ajax我的实际应用例子

复制代码
1 //1.$.ajax带json数据的异步请求
2 var aj = $.ajax( {  
		3     url:'productManager_reverseUpdate',// 跳转到 action  
		4     data:{  
		5              selRollBack : selRollBack,  
		6              selOperatorsCode : selOperatorsCode,  
		7              PROVINCECODE : PROVINCECODE,  
		8              pass2 : pass2  
		9     },  
		10     type:'post',  
		11     cache:false,  
		12     dataType:'json',  
		13     success:function(data) {  
		14         if(data.msg =="true" ){  
		15             // view("修改成功！");  
		16             alert("修改成功！");  
		17             window.location.reload();  
		18         }else{  
		19             view(data.msg);  
		20         }  
		21      },  
		22      error : function() {  
		23           // view("异常！");  
		24           alert("异常！");  
		25      }  
		26 });
		27 
		28 
		29 //2.$.ajax序列化表格内容为字符串的异步请求
		30 function noTips(){  
		31     var formParam = $("#form1").serialize();//序列化表格内容为字符串  
		32     $.ajax({  
		33         type:'post',      
		34         url:'Notice_noTipsNotice',  
		35         data:formParam,  
		36         cache:false,  
		37         dataType:'json',  
		38         success:function(data){  
		39         }  
		40     });  
		41 }  
		42 
		43 
		44 //3.$.ajax拼接url的异步请求
		45 var yz=$.ajax({  
		46      type:'post',  
		47      url:'validatePwd2_checkPwd2?password2='+password2,  
		48      data:{},  
		49      cache:false,  
		50      dataType:'json',  
		51      success:function(data){  
		52           if( data.msg =="false" ) //服务器返回false，就将validatePassword2的值改为pwd2Error，这是异步，需要考虑返回时间  
		53           {  
		54                textPassword2.html("<font color='red'>业务密码不正确！</font>");  
		55                $("#validatePassword2").val("pwd2Error");  
		56                checkPassword2 = false;  
		57                return;  
		58            }  
		59       },  
		60       error:function(){}  
		61 }); 
		62 
		63 
		64 //4.$.ajax拼接data的异步请求
		65 $.ajax({   
		66     url:'<%=request.getContextPath()%>/kc/kc_checkMerNameUnique.action',   
		67     type:'post',   
		68     data:'merName='+values,   
		69     async : false, //默认为true 异步   
		70     error:function(){   
		71        alert('error');   
		72     },   
		73     success:function(data){   
		74        $("#"+divs).html(data);   
		75     }
		76 });
