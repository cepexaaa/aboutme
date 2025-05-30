package ru.itmo.wp.controller;

import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.TicTacToe;
import ru.itmo.wp.service.TicTacToeService;

@RestController
@RequestMapping("/api/game")
public class TicTacToeController {
    private final TicTacToeService TicTacToeService;

    public TicTacToeController(TicTacToeService TicTacToeService) {
        this.TicTacToeService = TicTacToeService;
    }

    @PostMapping("/start")
    public TicTacToe startTicTacToe() {
        return TicTacToeService.startNewGame();
    }

    @PostMapping("/move")
    public TicTacToe makeMove(@RequestParam String gameId, @RequestParam int row, @RequestParam int col) {
        return TicTacToeService.makeMove(gameId, row, col);
    }

    @GetMapping("/status")
    public TicTacToe getGameStatus(@RequestParam String gameId) {
        return TicTacToeService.getGameStatus(gameId);
    }
}