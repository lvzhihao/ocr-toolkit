//demo.js
//
var uploadImage;
$( function () {
    uploadImage = $( '#uploadImage' ).imageEditer( {
        backImgUrl: '/public/h5-imgcrop/img/card_3.png'
    } );
    $( '#upload' ).on( 'click', function () {
        uploadImage.imageEditer( 'upload' )
    } )
    $( '#getimg' ).on( 'click', function () {
        var attr = uploadImage.imageEditer( 'attr' );
        var data = uploadImage.imageEditer( 'getImage' );
        $( '#preview' ).width( attr.map_w ).attr( {
            'src': data.imgurl
        } );
        console.log(attr)
        $( '#log' ).html( JSON.stringify( data.attr ) );
        /*
        x = attr.map_w * 0.3 * -1
        width = attr.map_w * 0.6
        y = $( '#preview' ).height() * 0.84 * -1
        height = $( '#preview' ).height() * 0.13
        */
        x = data.attr.canvas_width * 0.32
        y = data.attr.canvas_height * 0.78 
        width = data.attr.canvas_width * 0.62
        height = data.attr.canvas_height * 0.13
        if (x + width > data.attr.canvas_width ) {
            width = data.attr.canvas_width -x
        }
        if (y + height > data.attr.canvas_height ) {
            height = data.attr.canvas_width -y
        }


        console.log(x, y, width, height)

        var canvas = document.getElementById("canvas")
        canvas.width = width 
        canvas.height = height 
        ctx = canvas.getContext('2d')
        /*
        ctx.globalAlpha = 1.0
        ctx.fillStyle = 'white'
        ctx.fillRect(0, 0, width, height)
        ctx.save()
        */
        ctx.drawImage(document.getElementById('preview'), x, y, width, height, 0, 0, width, height );

        trimCanvas = trim(canvas)
        //console.log(canvas.width, canvas.height, x / data.attr.scale, y / data.attr.scale)

        imgdata = canvas.toDataURL()

        $.ajax({
            url: "/api/demo",
            method: "POST",
            contentType: "application/json; charset=UTF-8",
            data: JSON.stringify({
                "image":  imgdata.replace(/.*;base64,(.*)/mg, "\$1"),
                "whitelist": "0123456789xX"
            })
        }).done(function(msg){
            console.log(msg)
            alert(msg)
        });
    } )
} )
//旋转	
function rotateImage( val ) {
    console.log( val )
        $( '#rotatevalue' ).html( val );
    $( '#result_img' ).css( {
        'transform': 'scale(' + $( '#scalevalue' ).html() + ') rotate(' + $( '#rotatevalue' ).html() +
                'deg)'
    } )
}
//缩放
function scaleImage( val ) {
    console.log( val )
        $( '#scalevalue' ).html( val );
    $( '#result_img' ).css( {
        'transform': 'scale(' + $( '#scalevalue' ).html() + ') rotate(' + $( '#rotatevalue' ).html() +
                'deg)'
    } )
}
//水平
function vImage( val ) {
    $( '#vvalue' ).html( val );
    $( '#result_img' ).css( {
        'top': ( 1.92188 + parseFloat( $( '#hvalue' ).html() ) ) + 'px',
        'left': ( -199.352 + parseFloat( $( '#vvalue' ).html() ) ) + 'px'
    } )
}
//垂直
function hImage( val ) {
    $( '#hvalue' ).html( val );
    $( '#result_img' ).css( {
        'top': ( 1.92188 + parseFloat( $( '#hvalue' ).html() ) ) + 'px',
        'left': ( -199.352 + parseFloat( $( '#vvalue' ).html() ) ) + 'px'
    } )
}

function trim(c) {
    var ctx = c.getContext('2d'),
    copy = document.createElement('canvas').getContext('2d'),
    pixels = ctx.getImageData(0, 0, c.width, c.height),
    l = pixels.data.length,
    i,
    bound = {
        top: null,
        left: null,
        right: null,
        bottom: null
    },
    x, y;

    for (i = 0; i < l; i += 4) {
        if (pixels.data[i+3] !== 0) {
            x = (i / 4) % c.width;
            y = ~~((i / 4) / c.width);

            if (bound.top === null) {
                bound.top = y;
            }

            if (bound.left === null) {
                bound.left = x; 
            } else if (x < bound.left) {
                bound.left = x;
            }

            if (bound.right === null) {
                bound.right = x; 
            } else if (bound.right < x) {
                bound.right = x;
            }

            if (bound.bottom === null) {
                bound.bottom = y;
            } else if (bound.bottom < y) {
                bound.bottom = y;
            }
        }
    }

    var trimHeight = bound.bottom - bound.top,
    trimWidth = bound.right - bound.left,
    trimmed = ctx.getImageData(bound.left, bound.top, trimWidth, trimHeight);

    copy.canvas.width = trimWidth;
    copy.canvas.height = trimHeight;
    copy.putImageData(trimmed, 0, 0);

    // open new window with trimmed image:
    return copy.canvas;
}
