String.prototype.capitalize = function() {
	return this.charAt(0).toUpperCase() + this.slice(1);
}

!function( $ ) {
	$(function () {
		"use strict"
		var getContents = function(fracType) {
			$('#config').load('/' + fracType + ' form', function() {
				var form = $('#config form');

				$('#maps').width($(window).width())
				          .height($(window).height());
				var map = initialize(fracType);
				$('input', form).change(function() {
					console.log("Form changed, redrawing");
					var z = map.getZoom();
					map.setZoom(z + 1);
					map.setZoom(z);
				});
				$(window).resize(function() {
					google.maps.event.trigger(map, 'resize')
				});
			});
		};

		$('ul.nav li a').click(function(e) {
			var fracType = $(this).attr('id');
			console.log(fracType);
			getContents(fracType);
			$('#masthead').fadeOut();
			$('#config').fadeIn();
			return false;
		});
	})
}( window.jQuery );
