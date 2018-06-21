var $container = $('#image-container');
var $status = $('#status');
var $progress = $('progress');

var supportsProgress = $progress[0] &&

	$progress[0].toString().indexOf('Unknown') === -1;

var loadedImageCount, imageCount;

$('#add').click( function() {

	var items = getItems();
	$container.prepend( $(items) );

	$container.imagesLoaded()
		.progress( onProgress )
		.always( onAlways );

	imageCount = $container.find('img').length;
	resetProgress();
	updateProgress( 0 );
});


$('#reset').click( function() {
	$container.empty();
});


function getItems() {
	var items = '';
	for ( var i = 0; i < 7; i++ ) {
		items += getImageItem();
	}
	return items;
}

function getImageItem() {
	var item = '<li>';
	$.ajax({
		type:"get",
		url:"getoneimage",
		async:false,
		data: {
			name : "liwei",
			passwrod : "password",
		},
		success:function(data) {
			item += '<img src="' + data + '" /></li>';
		}
	});

	return item;
}

function resetProgress() {
	$status.css({ opacity: 1 });
	loadedImageCount = 0;
	if ( supportsProgress ) {
		$progress.attr( 'max', imageCount );
	}
}

function updateProgress( value ) {
	if ( supportsProgress ) {
		$progress.attr( 'value', value );
	} else {
		$status.text( value + ' / ' + imageCount );
	}
}

function onProgress( imgLoad, image ) {
	var $item = $( image.img ).parent();
	$item.removeClass('is-loading');
	if ( !image.isLoaded ) {
		$item.addClass('is-broken');
	}

	loadedImageCount++;
	updateProgress( loadedImageCount );
}

function onAlways() {
	$status.css({ opacity: 0 });
}

