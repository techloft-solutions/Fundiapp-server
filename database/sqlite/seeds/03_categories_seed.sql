INSERT INTO `categories` (
    `name`,
    `parent_id`,
    `icon_url`,
    `level`,
    `industry_id`
) VALUES (
    'Artisan',
    NULL,
    'https://www.clipartmax.com/png/middle/106-1061107_plumber-free-icon-construction.png',
    0,
    NULL
), (
    'Professional',
    NULL,
    'https://www.clipartmax.com/png/middle/106-1061107_plumber-free-icon-construction.png',
    0,
    NULL
), (
    'Optician',
    1,
    'https://www.clipartmax.com/png/middle/106-1061107_plumber-free-icon-construction.png',
    1,
    1
), (
    'Surgeon',
    2,
    'https://www.clipartmax.com/png/middle/106-1061107_plumber-free-icon-construction.png',
    1,
    1
), (
    'Nurse',
    1,
    'https://www.clipartmax.com/png/middle/106-1061107_plumber-free-icon-construction.png',
    1,
    1
);