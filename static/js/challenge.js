$(".chall-preview").click(function(e) {
    var modal = $("div." + e.currentTarget.name);
    modal.removeClass("d-none");
    modal.find(".chall-netcat-copy").removeClass("bi-clipboard-check-fill").addClass("bi-clipboard-fill");
    var contents = modal.contents().find(".challenge-flag");
    if(contents.length > 0) {
        contents[0].focus();
    }
    const overlay = document.getElementById('overlay');
    overlay.style.display = 'block';

    e.preventDefault();
});

$(document).mouseup(function(e) {
    var container = $("div.chall-toggle");
    const overlay = document.getElementById('overlay');
    if (!container.is(e.target) && container.has(e.target).length === 0) {
        container.addClass("d-none");
        overlay.style.display = 'none';
    }
});

$(".chall-close").click(function(e) {
    $(e.currentTarget).parent().addClass("d-none");
    const overlay = document.getElementById('overlay');
    overlay.style.display = 'none';
});

$(".chall-netcat-copy").click(function(e) {
    var text = $(e.currentTarget).text();
    console.log(text);
    console.log("copied.");
    navigator.clipboard.writeText(text).then(function() {
        $(e.currentTarget).removeClass("bi-clipboard-fill");
        $(e.currentTarget).addClass("bi-clipboard-check-fill");
    }, function(err) {
        console.error('Could not copy text: ', err);
    });
});

$(".submit-form").submit(async function(e){
    e.preventDefault();

    e = $(e.currentTarget).find(".chall-submit");


    var correct = false;

    var flag = e.parents(".input-group-append").siblings(".chall-flag").val();
    var chall_id = e.attr('name');

    var formData = new FormData();
    formData.append("flag", flag);
    formData.append("challID", chall_id);
    
    var response = await fetch("/submit", {
        method: "POST",
        body: formData
    });
    

    if (response.status == 202) {
        //flag correct
        correct = true;
    }
    else if (response.status == 409) {
        correct = false;
    }
    else if (response.status == 406) {
        correct = false;
    }
    else {
        console.error("Error:", response);
    }

    var chall_name = e.parents().parents().siblings(".chall-dial-name").text().trim();

    if (correct){
        submitFlag(e.currentTarget);
        e.parents(".chall-toggle").addClass("chall-toggle-solved");
        e.prop('disabled', true);
        e.parents(".input-group-append").siblings(".chall-flag").prop('disabled', true);
        var button = $(".chall-preview").filter(function() {
            return $(this).find(".chall-title").text().trim() === chall_name;
        });
        button.addClass("chall-solved");
    }
    else {
        e.parents(".chall-toggle").addClass("chall-toggle-wrong");
        e.parents(".input-group-append").siblings(".chall-flag")
            .addClass("text-danger")
            .css('outline', '1px solid red');
        e.addClass("btn-danger");
        setTimeout(function(){
            e.parents(".chall-toggle").removeClass("chall-toggle-wrong");
            e.parents(".input-group-append").siblings(".chall-flag")
                .removeClass("text-danger")
                .css('outline', 'none');
            e.removeClass("btn-danger");
        },1000);
    }
});

function submitFlag(form){
   challengeSolved();
}

function startFireworks(){
    const container = document.querySelector('.overlay');
    const fireworks = new Fireworks.default(container, {
        autoresize: true,
        opacity: 1,
        acceleration: 1.05,
        friction: 0.97,
        gravity: 1.5,
        particles: 50,
        trace: 3,
        explosion: 5,
        intensity: 30,
        flickering: 50,
        lineStyle: 'round',
        hue: {
            min: 0,
            max: 360
        },
        delay: {
            min: 30,
            max: 60
        },
        rocketsPoint: {
            min: 50,
            max: 50
        },
        lineWidth: {
            explosion: {
                min: 1,
                max: 3
            },
            trace: {
                min: 1,
                max: 2
            }
        },
        brightness: {
            min: 50,
            max: 80
        },
        decay: {
            min: 0.015,
            max: 0.03
        },
        sound: {
            enable: true,
            files: [
                'https://fireworks.js.org/sounds/explosion0.mp3',
                'https://fireworks.js.org/sounds/explosion1.mp3',
                'https://fireworks.js.org/sounds/explosion2.mp3'
            ],
            volume: {
                min: 4,
                max: 8
            }
        }
    });
    const intervals = [100, 200, 200, 300, 300, 300];
    intervals.forEach((interval, index) => {
        setTimeout(() => {
            fireworks.launch(1);
        }, interval * index);
    });
    setTimeout(() => {
        fireworks.clear();
        while (container.firstChild) {
            container.removeChild(container.firstChild);
        }
    }, 5000);
}

function challengeSolved(challId){
    startFireworks();
}

document.addEventListener('DOMContentLoaded', function() {
    const overlay = document.getElementById('overlay');
    const challengeToggles = document.querySelectorAll('.chall-toggle');
    const closeButtons = document.querySelectorAll('.chall-close');

    challengeToggles.forEach(toggle => {
        toggle.addEventListener('show.bs.collapse', function() {
            overlay.style.display = 'block';
        });
    });

    closeButtons.forEach(button => {
        button.addEventListener('click', function() {
            overlay.style.display = 'none';
        });
    });
});