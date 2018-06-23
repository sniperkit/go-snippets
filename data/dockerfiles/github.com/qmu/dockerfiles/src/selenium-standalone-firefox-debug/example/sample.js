
"use strict";

const webdriver = require('selenium-webdriver'),
    expect = require('expect.js'),
    assert = require('assert'),
    fs = require('fs'),
    By = webdriver.By,
    until = webdriver.until,
    t = require('selenium-webdriver/testing'),
    firefox = require('selenium-webdriver/firefox');

var url = "http://selenium-standalone-firefox-debug:4444/wd/hub";

t.describe('Selenium Test Demo', () => {
    let driver;

    t.before(() => {
        driver = new webdriver.Builder()
            .usingServer(url)
            .withCapabilities(webdriver.Capabilities.firefox())
            .build();
        let w = driver.manage().window();
        driver.manage().timeouts().implicitlyWait(500000);
        w.setSize(375, 667);
    });

    t.after(() => {
        driver.quit();
    });

    t.it('Take a screenshot of google', () => {
        driver.get('https://www.google.co.jp').then(() => {
            driver.takeScreenshot().then((image, err) => {
                fs.writeFile('./output/out.png', image, 'base64', (err) => {
                    console.log(err);
                });
            });
        });
    });
});
