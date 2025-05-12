#coding="utf-8"

import asyncio
from pyppeteer import launch

async def main():
    browser = await launch({"headless": True})
    page = await browser.newPage()
    await page.goto('https://www.linkedin.com/signup')
    email_input=await page.querySelector("input[name=email-address]")#获取页面上某个元素
    pass_input=await page.querySelector("input#password")
    agree_button=await page.querySelector("button#join-form-submit")
    await email_input.type("123456@qq.com")  #在input里输入内容
    await pass_input.type("123456") #点击
    text=await agree_button.getProperty("id")
    print(text)
    await agree_button.click()


    # html=await page.content() #查看整个页面的html
    # print(html)
    await page.waitForSelector("input#last-name",{'visible': True}) #等页面上出现某个特定元素
    first_name_input=await page.querySelector("input#first-name")
    last_name_input=await page.querySelector("input#last-name")
    continue_button=await page.querySelector("button#join-form-submit")
    await first_name_input.type("dmk")
    await last_name_input.type("laos")
    text=await continue_button.getProperty("id")
    print(text)
    await continue_button.click()

    await page.waitFor(4000)        #休息4000ms
    await page.screenshot({'path': 'screenshot.png'})   #截屏保存
    # await browser.close()

asyncio.get_event_loop().run_until_complete(main())

