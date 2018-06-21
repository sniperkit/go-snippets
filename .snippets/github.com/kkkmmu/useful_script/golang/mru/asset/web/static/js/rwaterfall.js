window.addEventListener('DOMContentLoaded', function(e) {

    var waterfall = new Waterfall({ minBoxWidth: 200 });
    window.onscroll = function() {
        var i = waterfall.getHighestIndex();
        if(i > -1) {
            // get last box of the column
            var lastBox = Array.prototype.slice.call(waterfall.columns[i].children, -1)[0];
            if(checkSlide(lastBox)) {
                var count = 5;
                while(count--) waterfall.addBox(boxHandle());
            }
        }
    };

    function checkSlide(elem) {
        if(elem) {
            var screenHeight = (document.documentElement.scrollTop || document.body.scrollTop) +
                               (document.documentElement.clientHeight || document.body.clientHeight);
            var elemHeight = elem.offsetTop + elem.offsetHeight / 2;

            return elemHeight < screenHeight;
        }
    }

    function newNode() {
        var size = ['660x250', '300x400', '350x500', '200x320', '300x300'];
        var color = [ 'E97452', '4C6EB4', '449F93', 'D25064', 'E59649' ];
        var i = 0;

        return function() {
			var imgsrc ="" ;
            var box = document.createElement('div');
            box.className = 'wf-box';
            var image = document.createElement('img');
            image.src = imgsrc;
            box.appendChild(image);
            var content = document.createElement('div');
            content.className = 'content';
            var title = document.createElement('p');
            title.appendChild(document.createTextNode('Heading'));
            content.appendChild(title);
            box.appendChild(content);

            if(++i === size.length) i = 0;

            return box;
        };
    }
});
