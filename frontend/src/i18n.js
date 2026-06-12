import { addMessages, init, locale } from 'svelte-i18n';
import en from './locales/en.json';
import ru from './locales/ru.json';

addMessages('en', en);
addMessages('ru', ru);

const savedLocale = localStorage.getItem('locale') || 'ru';

init({
  fallbackLocale: 'en',
  initialLocale: savedLocale,
});

locale.subscribe((value) => {
  if (value) {
    localStorage.setItem('locale', value);
  }
});
