sdes.utils.date = function() {
    "use strict";

    var _this   = this,
        msInDay = 1000 * 60 * 60 * 24;

    this.timeToYearMonthDay = function(time, padd) {
        var date  = time === undefined || time === null ? new Date() : new Date(time),
            year  = date.getFullYear(),
            month = date.getMonth()+1,
            day   = date.getDate();

        if ( padd ) {
            if ( month < 10 )
                month = "0"+month;

            if ( day < 10 )
                day = "0"+day;
        }

        return year+"-"+month+"-"+day;
    }

    this.getDaysBetweenDates = function(from, to)  {
        var utc1 = Date.UTC(from.getFullYear(), from.getMonth(), from.getDate());
        var utc2 = Date.UTC(to.getFullYear(), to.getMonth(), to.getDate());

        return Math.floor((utc2 - utc1)/msInDay);
    }

    this.addDaysToDate = function(date, days) {
        var d = new Date(date);
        d.setDate(d.getDate()+days);
        return d;
    }
}
