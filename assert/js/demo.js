//demo.js
//
var uploadImage;
$( function () {
    uploadImage = $( '#uploadImage' ).imageEditer( {
        backImgUrl: '/public/h5-imgcrop/img/card_1.png'
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
        $( '#log' ).html( JSON.stringify( data.attr ) );
        x = attr.map_w * 0.3 * -1
        width = attr.map_w * 0.6
        y = $( '#preview' ).height() * 0.84 * -1
        height = $( '#preview' ).height() * 0.13
        console.log(x, y, width, height)

        var canvas = document.getElementById("canvas")
        canvas.width = width / data.attr.scale
        canvas.height = height / data.attr.scale
        ctx = canvas.getContext('2d')
        ctx.drawImage(document.getElementById('preview'), x / data.attr.scale, y / data.attr.scale);
        console.log(x / data.attr.scale, y / data.attr.scale)

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

