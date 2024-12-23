import 'htmx.org';
import Alpine from 'alpinejs';
import NProgress from 'nprogress';
import Toastify from 'toastify-js';

// init alpine
window.Alpine = Alpine
Alpine.start()

// show progress bar on htmx request
NProgress.configure({ showSpinner: false });
document.addEventListener('htmx:beforeSend', function (e) {
  if (NProgress.isStarted()) return;
  NProgress.start();
});
document.addEventListener('htmx:afterRequest', function (e) {
  // If the response has a redirect header, keep the progress bar running (until the redirect is complete)
  if (e.detail.xhr.getResponseHeader('HX-Location')) return;
  NProgress.done();
});

// handle errors
document.addEventListener('htmx:sendError', function (e) {
  // Network errors don't have content to swap, so display a toast instead
  Toastify({
    text: 'A network error occurred',
    duration: 0,
    newWindow: true,
    close: true,
    gravity: 'top',
    position: 'right',
    stopOnFocus: true,
    style: {
      background: '#f97316',
    },
  }).showToast();
});
document.addEventListener('htmx:beforeSwap', function (e) {
  // Responses with non-200 status codes are usually in JSON format, not HTML,
  // so we don't want to swap them with something else.
  const status = e.detail.xhr.status;
  if (status > 200) {
    // Instead of swapping, we'll show a toast
    e.detail.shouldSwap = false;
    const body = JSON.parse(e.detail.xhr.responseText);
    if (body.error) {
      Toastify({
        text: body.error,
        duration: 0,
        newWindow: true,
        close: true,
        gravity: 'top',
        position: 'right',
        stopOnFocus: true,
        style: {
          background: '#f97316',
        },
      }).showToast();
    }
  }
});
