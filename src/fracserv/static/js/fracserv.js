String.prototype.capitalize = function() {
	return this.charAt(0).toUpperCase() + this.slice(1);
}

!function( $ ) {
	$(function () {
		"use strict"
		var getContents = function(fracType) {
			$('#config-content').load('/' + fracType + ' form', function() {
				var form = $('#config-content form');
				var map = initialize(fracType);
				var emptyForm = $('input', form).length == 0;
				$('#config').toggle(!emptyForm);
				if(!emptyForm) {
					$('input', form).bind('input', function() {
						console.log("Form changed, redrawing");
						map.reload();
					});

					form.submit(function() {
						console.log("Form submitted");
						map.reload();
						return false;
					});
				}

				var resize = function() {
					console.log("resize");
					var navHeight = $('.navbar-fixed-top').height();
					$('#maps').width($(window).width())
							  .height($(window).height()-navHeight)
							  .css('top', navHeight+'px');
					google.maps.event.trigger(map, 'resize')
				}
				$(window).resize(function() {
					resize();
				});
				$('a.btn.btn-navbar').click(resize);
				resize();
			});
		};
		var dismiss = function() {
			$('#config').fadeOut();
			$('#gear').fadeIn();
		};

		var show = function() {
			$('#config').fadeIn();
			$('#gear').fadeOut();
		};
		$('#hide').click(dismiss);
		$('#gear').click(show);

		$('ul.nav li a').click(function(e) {
			var fracType = $(this).attr('id');
			console.log(fracType);
			getContents(fracType);
			$('#masthead').fadeOut();
			$('#mobile-jump').fadeOut();
			$('#config').fadeIn();
			return false;
		});
	})
}( window.jQuery );
