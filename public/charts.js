function i2h(i) {
    var h = parseInt(i).toString(16);
    return (h.length < 2) ? "0" + h : h;
}

function makePlotRange(min, max, count) {
    var num = Math.floor(max / count);
    var i, l;
    var a = [];

    for (i = 0, l = count; i < l; i++) {
	var r, g, inc = 127;


	if (min >= (max / 2)) {
	    r = Math.floor((255 * min) / inc);
	    g = Math.floor((255 * (inc - min)) / inc);
	    if (g < 0) {
		g = 0;
	    }
	    if (r > 255) {
		r = 255;
	    }
	} else {
	    g = Math.floor((255 * min) / inc);
	    r = Math.floor((255 * (inc - min)) / inc);
	    if (r < 0) {
		r = 0;
	    }
	    if (g > 255) {
		g = 255;
	    }
	}

	a.push({from: min, to: min + num, color: '#' + i2h(r) + i2h(g) + '00'});
	min += num;
    }
    
    return a;
}


function makeSpeed(title, min, max, count) {
    return {
	    chart: {
		type: 'gauge',
		plotBackgroundColor: null,
		plotBackgroundImage: null,
		plotBorderWidth: 0,
		plotShadow: false
	    },

	    title: {
		text: title
	    },

	    pane: {
		startAngle: -150,
		endAngle: 150,
		background: [{
		    backgroundColor: {
			linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
			stops: [
			    [0, '#FFF'],
			    [1, '#333']
			]
		    },
		    borderWidth: 0,
		    outerRadius: '109%'
		}, {
		    backgroundColor: {
			linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
			stops: [
			    [0, '#333'],
			    [1, '#FFF']
			]
		    },
		    borderWidth: 1,
		    outerRadius: '107%'
		}, {
		    // default background
		}, {
		    backgroundColor: '#DDD',
		    borderWidth: 0,
		    outerRadius: '105%',
		    innerRadius: '103%'
		}]
	    },

	    // the value axis
	    yAxis: {
		min: 20,
		max: 120,

		minorTickInterval: 'auto',
		minorTickWidth: 1,
		minorTickLength: 10,
		minorTickPosition: 'inside',
		minorTickColor: '#666',

		tickPixelInterval: 30,
		tickWidth: 2,
		tickPosition: 'inside',
		tickLength: 10,
		tickColor: '#666',
		labels: {
		    step: 2,
		    rotation: 'auto'
		},
		title: {
		    text: 'F°'
		},
		plotBands: makePlotRange(min, max, count)
	    },
	
	    series: [{
		name: title,
		data: [min],
		tooltip: {
		    valueSuffix: ' °F'
		}
	    }]
    };
}

function makeChart(name, attr) {
    return function() {
	$('#gauge_'+name).highcharts(
	    makeSpeed(name, 0, 120, 6),
	    function(chart) {
		if (!chart.renderer.forExport) {
		    setInterval(function() {
			$.get({
			    url: 'data/current/' + name,
			    success: function(data) {
				var point = chart.series[0].points[0];
				point.update(Math.floor(data[attr]));
			    }
			});
		    }, 10000);
		}
	    });
    };
}

$(function () {
    $.get({
	url: 'data/sensors',
	success: function(data) {
	    var i, l = data.length ;
	    for (i = 0; i < l; i++) {
		console.log(data[i].name);
		makeChart(data[i].name, "temp")();
	    }
	}
    });
});


$(function () {
    $('#linedata').highcharts({
	chart: {
	    type: 'spline'
	},
	title: {
	    text: 'Weather data for Greenhouse'
	},
	subtitle: {
	    text: 'May 31 and and June 1, 2015 at two locations in Vik i Sogn, Norway'
	},
	xAxis: {
	    type: 'datetime',
	    labels: {
		overflow: 'justify'
	    }
	},
	yAxis: {
	    title: {
		text: 'Ideal Conditions'
	    },
	    minorGridLineWidth: 0,
	    gridLineWidth: 0,
	    alternateGridColor: null,
	    plotBands: [{ // Light air
		from: 0.3,
		to: 1.5,
		color: 'rgba(68, 170, 213, 0.1)',
		label: {
		    text: 'Too Cold',
		    style: {
			color: '#606060'
		    }
		}
	    }, { // Light breeze
		from: 1.5,
		to: 3.3,
		color: 'rgba(0, 0, 0, 0)',
		label: {
		    text: 'Cold',
		    style: {
			color: '#606060'
		    }
		}
	    }, { // Gentle breeze
		from: 3.3,
		to: 5.5,
		color: 'rgba(68, 170, 213, 0.1)',
		label: {
		    text: 'Just Right',
		    style: {
			color: '#606060'
		    }
		}
	    }, { // Moderate breeze
		from: 5.5,
		to: 8,
		color: 'rgba(0, 0, 0, 0)',
		label: {
		    text: 'Hot',
		    style: {
			color: '#606060'
		    }
		}
	    }, { // Fresh breeze
		from: 8,
		to: 11,
		color: 'rgba(68, 170, 213, 0.1)',
		label: {
		    text: 'Too Hot',
		    style: {
			color: '#606060'
		    }
		}
	    }]
	},
	tooltip: {
	    valueSuffix: ' °F'
	},
	plotOptions: {
	    spline: {
		lineWidth: 4,
		states: {
		    hover: {
			lineWidth: 5
		    }
		},
		marker: {
		    enabled: false
		},
		pointInterval: 3600000, // one hour
		pointStart: Date.UTC(2015, 4, 31, 0, 0, 0)
	    }
	},
	series: [{
	    name: 'Temperature',
	    data: [0.2, 0.8, 0.8, 0.8, 1, 1.3, 1.5, 2.9, 1.9, 2.6, 1.6, 3, 4, 3.6, 4.5, 4.2, 4.5, 4.5, 4, 3.1, 2.7, 4, 2.7, 2.3, 2.3, 4.1, 7.7, 7.1, 5.6, 6.1, 5.8, 8.6, 7.2, 9, 10.9, 11.5, 11.6, 11.1, 12, 12.3, 10.7, 9.4, 9.8, 9.6, 9.8, 9.5, 8.5, 7.4, 7.6]

	}, {
	    name: 'Humidity',
	    data: [0, 0, 0.6, 0.9, 0.8, 0.2, 0, 0, 0, 0.1, 0.6, 0.7, 0.8, 0.6, 0.2, 0, 0.1, 0.3, 0.3, 0, 0.1, 0, 0, 0, 0.2, 0.1, 0, 0.3, 0, 0.1, 0.2, 0.1, 0.3, 0.3, 0, 3.1, 3.1, 2.5, 1.5, 1.9, 2.1, 1, 2.3, 1.9, 1.2, 0.7, 1.3, 0.4, 0.3]
	}],
	navigation: {
	    menuItemStyle: {
		fontSize: '10px'
	    }
	}
    });
});

function lineSeries(data) {

    var o = {}, d, i, len = data.length, k, count = 0;
    for (i = 0; i < len; i++) {
	d = JSON.parse(data[i]);
	if (count === 0) {
	    for (k in d) {
		o[k] = {};
		o[k].data = [];
		o[k].name = k;
	    }
	}
	
	for (k in len) {
	    o[k].data.push(d[k]);
	    console.log(d[k]);
	}
	
	count++;
    }

    return o;
}


console.log(lineSeries(lineData));
