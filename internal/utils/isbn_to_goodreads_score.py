import sys
from time import sleep
from selenium import webdriver
from selenium.webdriver.common.by import By
from webdriver_manager.firefox import GeckoDriverManager

def main():
    goodreads_url = "https://goodreads.com"
    goodreads_search_bar_element_xpath = '//*[@id="sitesearch_field"]'
    goodreads_score_xpath = '/html/body/div[2]/div[3]/div[1]/div[2]/div[4]/div[1]/div[2]/div[2]/span[2]'
    goodreads_search_icon_xpath = '/html/body/div[2]/div[2]/div[1]/div[1]/div[2]/div[4]/div[1]/div/div[2]/form/a'
    isbn = ''
    try:
        isbn = sys.argv[1]
    except:
        print("you should pass ISBN as argument to program")
        exit(1)

    cfg = webdriver.FirefoxOptions()
    cfg.headless = False
    driver = webdriver.Firefox(executable_path=GeckoDriverManager().install(),options=cfg)
    driver.get(goodreads_url)

    search_bar_element = driver.find_element(By.XPATH,goodreads_search_bar_element_xpath)
    search_bar_element.send_keys(isbn)
    search_icon_element = driver.find_element(By.XPATH,goodreads_search_icon_xpath)
    driver.execute_script("arguments[0].click();", search_icon_element)
    while not is_element_exists(driver,goodreads_score_xpath):
        sleep(1)

    score_element = driver.find_element(By.XPATH,goodreads_score_xpath)
    
    print(score_element.text)
def is_element_exists(driver:webdriver.Firefox,path):
    try:
        e = driver.find_element(By.XPATH,path)
        return True
    except:
        return False


if __name__ == "__main__":
    main()