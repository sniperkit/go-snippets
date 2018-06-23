
"use strict";

const webdriver = require('selenium-webdriver'),
    expect = require('expect.js'),
    assert = require('assert'),
    fs = require('fs'),
    By = webdriver.By,
    until = webdriver.until,
    t = require('selenium-webdriver/testing'),
    chrome = require('selenium-webdriver/chrome');

var url = "http://selenium_hub:4444/wd/hub";

t.describe('Selenium Test Demo', () => {
    let driver;

    t.before(() => {
        let options = new chrome.Options();
        options.addArguments("--user-agent=Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A403 Safari/8536.25");
        driver = new webdriver.Builder()
            .usingServer(url)
            .withCapabilities(webdriver.Capabilities.chrome())
            .setChromeOptions(options)
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
