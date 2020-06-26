import config from "@/config/Config";

/**
 * This is a wrapper around the LocalStorage API. It adds a prefix in front
 * of every key and to be able to store objects and arrays it uses
 * JSON.stringify and JSON.parse.
 */
class Storage {
  static getItem(key) {
    return JSON.parse(localStorage.getItem(`${config.STORAGE_PREFIX}_${key}`));
  }
  static setItem(key, value) {
    return localStorage.setItem(`${config.STORAGE_PREFIX}_${key}`, JSON.stringify(value));
  }
};

export default Storage;
