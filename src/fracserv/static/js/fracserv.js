String.prototype.capitalize = function() {
	return this.charAt(0).toUpperCase() + this.slice(1);
}

!function( $ ) {
	$(function () {
		"use strict"
		var getContents = function(fracType) {
			$('#config').load('/' + fracType + ' form', function() {
				var form = $('#config form');

				var map = initialize(fracType);
				$('input', form).bind('input', function() {
					console.log("Form changed, redrawing");
					map.reload();
				});

				var resize = function() {
					$('#maps').width($(window).width())
							  .height($(window).height());
				}
				resize();
				$(window).resize(function() {
					resize();
					google.maps.event.trigger(map, 'resize')
				});
				form.submit(function() {
					console.log("Form submitted");
					map.reload();
					return false;
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
