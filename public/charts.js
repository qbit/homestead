function i2h(i) {
    var h = parseInt(i).toString(16);
    return (h.length < 2) ? "0" + h : h;
}

var OM = new Date();
OM = OM.setMonth(OM.getMonth() - 1);

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
	credits: {
            enabled: false
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
	    data: [gaugeData[sensor + "_avg"]],
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
			    url: '/data/current/' + name,
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
	url: '/data/sensors',
	success: function(data) {
	    var i, l = data.length ;
	    for (i = 0; i < l; i++) {
		if (data[i].name === sensor) {
		  console.log(data[i].name);
		  makeChart(data[i].name, "temp")();
		}
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
	    text: 'Weather data for ' . sensor
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
	credits: {
            enabled: false
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
		pointInterval: 300000,
		pointStart: OM
	    }
	},
	series: lineSeries(lineData),
	navigation: {
	    menuItemStyle: {
		fontSize: '10px'
	    }
	}
    });
});

var lsMap = {
    "temp": "Temperature",
    "humidity": "Humidity",
    "pres_inhg": "Inch of mercury"
};

function lineSeries(data, map) {
    var o = {}, d, i, len = data.length, k, count = 0;
    for (i = 0; i < len; i++) {
	d = JSON.parse(data[i]);
	if (count === 0) {
	    for (k in d) {
		if (lsMap[k]) {
		    o[k] = {};
		    o[k].data = [];
		    o[k].name = lsMap[k];
		}
	    }
	}

	for (k in d) {
	    if (lsMap[k]) {
		o[k].data.push(d[k]);
	    }
	}

	count++;
    }

    var r = [];
    for (k in o) {
	r.push(o[k]);
    }

    return r;
}

var gaugeOptions = {

    chart: {
        type: 'solidgauge'
    },

    title: null,

    pane: {
        center: ['50%', '85%'],
        size: '140%',
        startAngle: -90,
        endAngle: 90,
        background: {
            backgroundColor: (Highcharts.theme && Highcharts.theme.background2) || '#EEE',
            innerRadius: '60%',
            outerRadius: '100%',
            shape: 'arc'
        }
    },

    tooltip: {
        enabled: false
    },

    // the value axis
    yAxis: {
        stops: [
            [0.1, '#55BF3B'], // green
            [0.5, '#DDDF0D'], // yellow
            [0.9, '#DF5353'] // red
        ],
        lineWidth: 0,
        minorTickInterval: null,
        tickAmount: 2,
        title: {
            y: -70
        },
        labels: {
            y: 16
        }
    },

    plotOptions: {
        solidgauge: {
            dataLabels: {
                y: 5,
                borderWidth: 0,
                useHTML: true
            }
        }
    }
};

$('#high-temp').highcharts(Highcharts.merge(gaugeOptions, {
    yAxis: {
        min: -100,
        max: 200,
        title: {
            text: 'High Temp'
        }
    },

    credits: {
        enabled: false
    },

    series: [{
        name: 'High Temp',
	data: [gaugeData[sensor + "_max"]],
        dataLabels: {
            format: '<div style="text-align:center"><span style="font-size:25px;color:' +
                ((Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black') + '">{y}</span><br/>' +
                '<span style="font-size:12px;color:silver">°F</span></div>'
        },
        tooltip: {
            valueSuffix: ' °F'
        }
    }]

}));

// The RPM gauge
$('#low-temp').highcharts(Highcharts.merge(gaugeOptions, {
    yAxis: {
        min: -100,
        max: 200,
        title: {
            text: 'Low Temp'
        }
    },

    credits: {
	enabled: false
    },

    series: [{
        name: 'Low Temp',
	data: [gaugeData[sensor + "_min"]],
        dataLabels: {
	    format: '<div style="text-align:center"><span style="font-size:25px;color:' +
                ((Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black') + '">{y:.1f}</span><br/>' +
                '<span style="font-size:12px;color:silver">°F</span></div>'
        },
        tooltip: {
	    valueSuffix: '°F'
        }
    }]

}));

