package ru.itmo.wp.form.validator;

import org.springframework.stereotype.Component;
import org.springframework.validation.Errors;
import org.springframework.validation.Validator;
import ru.itmo.wp.form.CommentCredentials;
import ru.itmo.wp.service.CommentService;

@Component
public class CommentCredentialsAddValidator implements Validator {
    private final CommentService commentService;

    public CommentCredentialsAddValidator(CommentService commentService) {
        this.commentService = commentService;
    }

    public boolean supports(Class<?> clazz) {
        return CommentCredentials.class.equals(clazz);
    }

    public void validate(Object target, Errors errors) {
        if (!errors.hasErrors()) {
            CommentCredentials commentForm = (CommentCredentials) target;
            if (!commentService.isTextVacant(commentForm.getText())) {
                errors.rejectValue("text", "text.isEmpty", "text is empty");
            }
        }
    }
}