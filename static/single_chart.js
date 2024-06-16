class MetricsChart {
    constructor(name, containerId, xData, yData, height, displayMin, displayMax) {
        this.name = name;
        this.containerId = containerId;
        this.container = document.getElementById(containerId);
        this.xData = xData;
        this.yData = yData;
        this.displayMin = displayMin;
        this.displayMax = displayMax;
        this.topHeaderHeight = 20;
        this.leftScaleWidth = 80;
        this.bottomScaleheight = 40;
        this.height = height;
        this.yValuesPadding = 0.05;
        this.colorGrid = '#333333';
        this.colorScalesText = '#AAAAAA';
        this.colorSeries = '#00CCFF';
        this.colorBackground = '#111111'
        this.colorDayBorder = '#FFFFFF';
        this.widthDayBorder = 1.5;
        this.widthSeries = 1.5;
        this.fontName = 'Roboto Mono'
        this.drawHorScale = true;
        this.canvas = document.createElement('canvas');
        this.context = this.canvas.getContext('2d');
        this.canvas.height = this.height;
        this.container.appendChild(this.canvas);
        this.resize();
        window.addEventListener('resize', () => this.resize());
    }

    resize() {
        this.width = this.container.clientWidth;
        this.canvas.width = this.width;
        this.drawChart();
    }

    setDisplayMinMax(displayMin, displayMax) {
        this.displayMin = displayMin;
        this.displayMax = displayMax;
        this.drawChart();
    }

    setData(xData, yData) {
        this.xData = xData;
        this.yData = yData;
        this.drawChart();
    }

    getXByValue(xValue) {
        const rangeValues = this.displayMax - this.displayMin;
        const rangePixels = this.width - this.leftScaleWidth
        const pixelsPerValue = rangePixels / rangeValues;
        const result = (xValue - this.displayMin) * pixelsPerValue;
        return this.leftScaleWidth + result;
    }

    getYByValue(yValue) {
        let maxDataValue = Math.max(...this.yData);
        let minDataValue = Math.min(...this.yData);
        const range = maxDataValue - minDataValue;
        minDataValue = minDataValue - range * this.yValuesPadding
        maxDataValue = maxDataValue + range * this.yValuesPadding

        const rangeValues = maxDataValue - minDataValue;
        const rangePixels = this.height - this.bottomScaleheight - this.topHeaderHeight;
        const pixelsPerValue = rangePixels / rangeValues;
        const result = (yValue - minDataValue) * pixelsPerValue;
        return (this.height - this.bottomScaleheight - this.topHeaderHeight) - result + this.topHeaderHeight;
        return result;
    }

    fromUnixTime(unixTime) {
        const date = new Date(unixTime * 1000);
        const year = date.getUTCFullYear();
        const month = String(date.getUTCMonth() + 1).padStart(2, '0');
        const day = String(date.getUTCDate()).padStart(2, '0');
        const hours = String(date.getUTCHours()).padStart(2, '0');
        const minutes = String(date.getUTCMinutes()).padStart(2, '0');
        const seconds = String(date.getUTCSeconds()).padStart(2, '0');
        return `${hours}:${minutes}`;
    }

    fromUnixTimeDate(unixTime) {
        const date = new Date(unixTime * 1000);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    }

    truncateNumberString(input) {
        const isNegative = input.startsWith('-');
        const numberString = isNegative ? input.slice(1) : input;

        if (numberString.length <= 10) {
            return input;
        }

        const [integerPart, fractionalPart] = numberString.split('.');

        if (parseFloat(input) < 1) {
            if (fractionalPart && fractionalPart.length > 10) {
                return (isNegative ? '-' : '') + `0.${fractionalPart.slice(0, 4)}..${fractionalPart.slice(-2)}`;
            } else if (fractionalPart) {
                return input.slice(0, 10);
            }
        } else {
            if (integerPart.length > 10) {
                return (isNegative ? '-' : '') + `${integerPart.slice(0, 4)}..${integerPart.slice(-2)}`;
            } else {
                return (isNegative ? '-' : '') + `${integerPart}.${fractionalPart.slice(0, 10 - integerPart.length - 1)}`;
            }
        }

        return input.slice(0, 10);
    }

    roundToNearestStep(seconds) {
        const HOUR = 3600;
        const STEPS = [HOUR / 60, HOUR / 30, HOUR / 20, HOUR / 12, HOUR / 6, HOUR / 4, HOUR / 2, HOUR, 2 * HOUR, 3 * HOUR, 6 * HOUR, 12 * HOUR, 24 * HOUR, 48 * HOUR, 72 * HOUR, 96 * HOUR];

        function getNearestStep(value, steps) {
            let nearest = steps[0];
            let minDiff = Math.abs(value - nearest);

            for (let i = 1; i < steps.length; i++) {
                let diff = Math.abs(value - steps[i]);
                if (diff < minDiff) {
                    minDiff = diff;
                    nearest = steps[i];
                }
            }
            return nearest;
        }

        return getNearestStep(seconds, STEPS);
    }

    roundToStep(value) {
        const STEPS = [1, 2, 5, 10, 15, 20, 50, 100, 200, 500, 1000]; // Массив шагов

        function getNearestStep(value, steps) {
            let nearest = steps[0];
            let minDiff = Math.abs(value - nearest);

            for (let i = 1; i < steps.length; i++) {
                let diff = Math.abs(value - steps[i]);
                if (diff < minDiff) {
                    minDiff = diff;
                    nearest = steps[i];
                }
            }
            return nearest;
        }

        const order = Math.floor(Math.log10(Math.abs(value)));
        const scaledSteps = STEPS.map(step => step * Math.pow(10, order - (step >= 10 ? 1 : 0)));

        const nearestStep = getNearestStep(value, scaledSteps);

        const roundedDown = Math.floor(value / nearestStep) * nearestStep;
        const roundedUp = Math.ceil(value / nearestStep) * nearestStep;

        return roundedUp;
    }

    formatNumbers(numbers) {
        function roundToDecimalPlaces(value, decimalPlaces) {
            const factor = Math.pow(10, decimalPlaces);
            const roundedValue = Math.round(value * factor) / factor;
            return roundedValue.toFixed(decimalPlaces).replace(/\.?0+$/, '');
        }

        function getMinDecimalPlaces(numbers) {
            let minDecimalPlaces = 0;
            for (let i = 0; i < numbers.length - 1; i++) {
                let diff = Math.abs(numbers[i] - numbers[i + 1]);
                let decimalPlaces = 0;
                while (diff < 1 && decimalPlaces < 20) {
                    diff *= 10;
                    decimalPlaces++;
                }
                minDecimalPlaces = Math.max(minDecimalPlaces, decimalPlaces);
            }
            return minDecimalPlaces;
        }

        numbers.sort((a, b) => a - b);

        const minDecimalPlaces = getMinDecimalPlaces(numbers);

        const result = numbers.map(number => {
            if (minDecimalPlaces === 0) {
                return Math.round(number).toString();
            } else {
                return roundToDecimalPlaces(number, minDecimalPlaces);
            }
        });

        return result;
    }

    clearRectWithBackgroundColor(context, x, y, width, height, backgroundColor) {
        context.clearRect(x, y, width, height);
        context.fillStyle = backgroundColor;
        context.fillRect(x, y, width, height);
    }

    truncateString(input, size) {
        if (input.length > size) {
            return input.slice(0, size - 2) + '..';
        }
        return input;
    }

    drawChart() {
        const ctx = this.context;
        const width = this.canvas.width;
        const height = this.canvas.height;
        const xData = this.xData;
        const yData = this.yData;

        ctx.save();

        let maxDataValue = Math.max(...yData);
        let minDataValue = Math.min(...yData);
        const range = maxDataValue - minDataValue;
        minDataValue = minDataValue - range * this.yValuesPadding
        maxDataValue = maxDataValue + range * this.yValuesPadding

        ctx.font = '12px ' + this.fontName

        ctx.imageSmoothingEnabled = true;
        ctx.imageSmoothingQuality = 'high';

        this.clearRectWithBackgroundColor(ctx, 0, 0, width, height, this.colorBackground)

        ////////////////////////////////////////////////////
        // Y SCALE /////////////////////////////////////////
        const yLabelHeight = 40;
        const yLabelCount = this.height / yLabelHeight;
        let yScaleStepTemp = (maxDataValue - minDataValue) / yLabelCount;
        let yScaleStep = this.roundToStep(yScaleStepTemp);

        let labels = []
        for (
            let yScaleValue = minDataValue - (minDataValue % yScaleStep);
            yScaleValue < maxDataValue;
            yScaleValue += yScaleStep) {
            labels.push(yScaleValue);
        }

        let labelsStr = this.formatNumbers(labels);


        ctx.lineJoin = 'bevel';
        ctx.beginPath();
        ctx.moveTo(this.leftScaleWidth - 3, 0);
        ctx.lineTo(this.leftScaleWidth - 3, this.height);
        ctx.strokeStyle = this.colorGrid;
        ctx.stroke();

        let index = 0;
        for (
            let yScaleValue = minDataValue - (minDataValue % yScaleStep);
            yScaleValue < maxDataValue;
            yScaleValue += yScaleStep) {
            const y = this.getYByValue(yScaleValue);

            ctx.save();
            ctx.rect(this.leftScaleWidth, 0, this.width - this.leftScaleWidth, this.height - this.bottomScaleheight - this.topHeaderHeight);
            ctx.clip();

            ctx.beginPath();
            ctx.moveTo(this.leftScaleWidth, y);
            ctx.lineTo(this.width, y);
            ctx.strokeStyle = this.colorGrid;
            ctx.stroke();

            ctx.restore();

            ctx.save();
            ctx.rect(0, this.topHeaderHeight, this.leftScaleWidth, this.height - this.bottomScaleheight - this.topHeaderHeight);
            ctx.clip();
            ctx.fillStyle = this.colorScalesText;
            let str = labelsStr[index];
            if (str.length > 10) {
                ctx.font = '10px ' + this.fontName
            }
            if (str.length > 13) {
                ctx.font = '8px ' + this.fontName;
            }
            if (str.length > 16) {
                ctx.font = '6px ' + this.fontName;
            }

            const metrics = ctx.measureText(str);
            const textWidth = metrics.width;
            let strX = this.leftScaleWidth - textWidth - 6;
            ctx.fillText(str, strX, y + 3);
            ctx.restore();

            index++;
        }
        ////////////////////////////////////////////////////

        ////////////////////////////////////////////////////
        // X SCALE /////////////////////////////////////////

        const metricsXScaleLabel = ctx.measureText("00.00.0000");
        const metricsXScaleLabelWidth = metricsXScaleLabel.width;
        const xLabelWidth = metricsXScaleLabelWidth * 1;

        const xLabelCount = (this.width - this.leftScaleWidth) / xLabelWidth;
        let xScaleStep = (this.displayMax - this.displayMin) / xLabelCount;
        xScaleStep = this.roundToNearestStep(xScaleStep);


        ctx.beginPath();
        ctx.moveTo(0, this.height - this.bottomScaleheight + 3);
        ctx.lineTo(this.width, this.height - this.bottomScaleheight + 3);
        ctx.strokeStyle = this.colorGrid;
        ctx.stroke();

        let dates = [];

        ctx.save();
        ctx.rect(this.leftScaleWidth, 0, this.width - this.leftScaleWidth, this.height - this.bottomScaleheight);
        ctx.clip();
        for (
            let xScaleValue = this.displayMin - (this.displayMin % 86400);
            xScaleValue < this.displayMax;
            xScaleValue += 86400) {
            const x = this.getXByValue(xScaleValue);
            ctx.beginPath();
            ctx.moveTo(x, this.height - this.bottomScaleheight);
            ctx.lineTo(x, 0);
            ctx.lineWidth = this.widthDayBorder;
            ctx.strokeStyle = this.colorDayBorder;
            ctx.stroke();
            dates.push(xScaleValue);
        }
        ctx.restore();

           for (
                let xScaleValue = this.displayMin - (this.displayMin % xScaleStep);
                xScaleValue < this.displayMax;
                xScaleValue += xScaleStep) {

                const x = this.getXByValue(xScaleValue);

                ctx.save();
                ctx.rect(this.leftScaleWidth, 0, this.width - this.leftScaleWidth, this.height - this.bottomScaleheight);
                ctx.clip();

                ctx.beginPath();
                ctx.moveTo(x, this.height - this.bottomScaleheight);
                ctx.lineTo(x, 0);
                ctx.strokeStyle = this.colorGrid;
                ctx.stroke();

                ctx.restore();

                if (this.drawHorScale) {

                const timeLabel = this.fromUnixTime(xScaleValue);

                ctx.save();
                ctx.rect(this.leftScaleWidth, this.height - this.bottomScaleheight, this.width - this.leftScaleWidth, this.bottomScaleheight);
                ctx.clip();
                ctx.fillStyle = this.colorScalesText;
                ctx.fillText(timeLabel, x - 18, this.height - this.bottomScaleheight + 15);
                ctx.restore();
                }
            }

        // Dates blocks
        if (this.drawHorScale) {
            for (let i in dates) {
                const dateLabel = this.fromUnixTimeDate(dates[i]);
                let posXbegin = this.getXByValue(dates[i]);
                let posXend = this.getXByValue(dates[i] + 86400);
                ctx.save();
                ctx.rect(this.leftScaleWidth, this.height - this.bottomScaleheight, this.width - this.leftScaleWidth, this.bottomScaleheight);
                ctx.clip();

                let y = this.height - this.bottomScaleheight + 25;

                ctx.beginPath();
                ctx.moveTo(posXbegin + 2, y);
                ctx.lineTo(posXend - 2, y);
                ctx.strokeStyle = this.colorScalesText;
                ctx.stroke();

                ctx.beginPath();
                ctx.moveTo(posXbegin + 2, y);
                ctx.lineTo(posXbegin + 2, y + 10);
                ctx.strokeStyle = this.colorScalesText;
                ctx.stroke();

                ctx.beginPath();
                ctx.moveTo(posXend - 2, y);
                ctx.lineTo(posXend - 2, y + 10);
                ctx.strokeStyle = this.colorScalesText;
                ctx.stroke();

                let textPos = posXbegin + (posXend - posXbegin) / 2;
                if (posXbegin < this.leftScaleWidth) {
                    textPos = this.leftScaleWidth + (posXend - this.leftScaleWidth) / 2;
                }
                if (posXend > this.width) {
                    textPos = posXbegin + (this.width - posXbegin) / 2;
                }

                textPos = textPos - 30;

                ctx.fillStyle = this.colorScalesText;
                ctx.fillText(dateLabel, textPos, y + 12);

                ctx.restore();
            }
        }


        ////////////////////////////////////////////////////
        //   DATA  /////////////////////////////////////////
        ctx.save();
        ctx.rect(this.leftScaleWidth, 0, this.width - this.leftScaleWidth, this.height - this.bottomScaleheight);
        ctx.clip();
        ctx.beginPath();
        yData.forEach((value, index) => {
            const x = this.getXByValue(this.xData[index]);
            const y = this.getYByValue(this.yData[index]);
            if (index === 0) {
                ctx.moveTo(x, y);
            } else {
                ctx.lineTo(x, y);
            }
        });
        ctx.strokeStyle = this.colorSeries;
        ctx.lineWidth = this.widthSeries;
        ctx.stroke();
        ctx.restore();
        ////////////////////////////////////////////////////

        ////////////////////////////////////////////////////
        //   LOGO  /////////////////////////////////////////
        /*ctx.fillStyle = this.colorBackground;
        ctx.fillRect(0, 0, this.width, 20);
        ctx.fillStyle = this.colorSeries;
        ctx.font = '16px ' + this.fontName
        ctx.fillText(this.name, 5, 14);
        ctx.font = '16px ' + this.fontName
        ctx.fillStyle = "#555555";
        ctx.fillText("u00.io - Blockchain analytics", this.width - 280, 14);*/
        ////////////////////////////////////////////////////

        ctx.restore();
    }
}
