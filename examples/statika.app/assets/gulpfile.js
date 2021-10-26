const {src, dest, watch, series, parallel} = require('gulp');
const sourcemaps = require('gulp-sourcemaps');
const sass = require('gulp-sass');
const concat = require('gulp-concat');
const uglify = require('gulp-uglify');
const postcss = require('gulp-postcss');
const autoprefixer = require('autoprefixer');
const cssnano = require('cssnano');
var replace = require('gulp-replace');
var rename = require('gulp-rename');
var imagemin = require('gulp-imagemin');
const del = require("del");
var multiDest = require('gulp-multi-dest');

const jsDists = ['../src/static/assets/js'];
const cssDists = ['../src/static/assets/css'];
const imgDists = ['../src/static/assets/img'];
const fontsDists = ['../src/static/assets/fonts'];

const files = {
    scssPath: 'src/scss/**/*.scss',
    jsPath: 'src/js/**/*.js',
    imgPath: 'src/img/**/*',
    fontsPath: 'src/fonts/**/*',
};

function scssTask() {
    return src(files.scssPath)
        .pipe(sourcemaps.init())
        .pipe(sass())
        .pipe(postcss([autoprefixer(), cssnano()]))
        // .pipe(sourcemaps.write('.'))
        .pipe(rename('site.css'))
        .pipe(multiDest(cssDists));
}

function jsTask() {
    return src([
        files.jsPath
    ])
        .pipe(concat('site.js'))
        .pipe(uglify())
        .pipe(multiDest(jsDists));
}

function imgTask() {
    return src(files.imgPath)
        .pipe(imagemin())
        .pipe(multiDest(imgDists));
}

function fontsTask() {
    return src(files.fontsPath)
        .pipe(multiDest(fontsDists));
}

function cacheBustTask() {
    var cbString = new Date().getTime();
    return src(['index.html'])
        .pipe(replace(/cb=\d+/g, 'cb=' + cbString))
        .pipe(dest('.'));
}

// Cleaning is evil...we output to other src folders
// function cleanTask() {
//   return del([
//       jsDists[0], jsDists[1],
//       cssDists[0], cssDists[1],
//       imgDists[0], imgDists[1],
//       fontsDists[0], fontsDists[1] ], {force: true});
// }

function watchTask() {
    watch([files.scssPath, files.jsPath, files.imgPath],
        {interval: 1000, usePolling: true}, //Makes docker work
        series(
            // To slow: parallel(scssTask, jsTask, imgTask, fontsTask),
            parallel(scssTask, jsTask),
            // cacheBustTask
        )
    );
}

exports.default = series(
    // cleanTask,
    parallel(scssTask, jsTask, imgTask, fontsTask),
    watchTask
);

// exports.clean = series(
//     cleanTask
// );
