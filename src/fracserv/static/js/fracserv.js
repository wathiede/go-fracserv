String.prototype.capitalize = function() {
	return this.charAt(0).toUpperCase() + this.slice(1);
}

!function( $ ) {
	$(function () {
		"use strict"
		var pop = $('#config');
		var getContents = function(fracType) {
			$.get('/' + fracType, function(data) {
				pop.attr('data-content', data);
				pop.attr('title', fracType.capitalize());

				pop.popover('show');

				$('form #w').val($(window).width());
				$('form #h').val($(window).height());


				$('form').submit(function() {
					$('body').css('background-image', 'url(/' + fracType + '?' + $(this).serialize() + ')');
					return false;
				});
				return data;
			});
		};

		pop.popover({
			placement: 'top',
			trigger: 'manual',
		});

		pop.click(function() {
			pop.popover('toggle');
		});

		$('ul.nav li a').click(function(e) {
			var fracType = $(this).attr('id');
			console.log(fracType);
			getContents(fracType);
			$('#masthead').fadeOut();
			return false;
		});
	})
}( window.jQuery );
