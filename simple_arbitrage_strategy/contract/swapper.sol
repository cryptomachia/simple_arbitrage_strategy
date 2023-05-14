pragma solidity ^0.8.0;

// SPDX-License-Identifier: MIT

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

interface ITraderJoe {
    function swapExactTokensForTokens(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external returns (uint256[] memory amounts);
}

interface ICurve {
    function exchange_underlying(
        int128 i,
        int128 j,
        uint256 dx,
        uint256 min_dy
    ) external;
}

interface IZyber {
    function swap(
        uint256 amount0Out,
        uint256 amount1Out,
        address to,
        bytes calldata data
    ) external;
}

contract Swaper {
    using SafeERC20 for IERC20;

    address private constant WETH = 0x82aF49447D8a07e3bd95BD0d56f35241523fBab1;
    address private constant USDT = 0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9;
    address private constant TRADERJOE_ROUTER = 0xd387c40a72703B38A5181573724bcaF2Ce6038a5;
    address private constant CURVE_POOL = 0x960ea3e3C7FB317332d990873d354E18d7645590;
    address private constant ZYBER_POOL = 0xC4e9FCe31518D4233C224772dC5532d49c5354c0;

    function swapTraderJoe(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        uint256 deadline
    ) external {
        IERC20(path[0]).safeTransferFrom(msg.sender, address(this), amountIn);
        IERC20(path[0]).safeApprove(TRADERJOE_ROUTER, amountIn);

        ITraderJoe(TRADERJOE_ROUTER).swapExactTokensForTokens(
            amountIn,
            amountOutMin,
            path,
            msg.sender,
            deadline
        );
    }

    function swapCurve(
        int128 i,
        int128 j,
        uint256 dx,
        uint256 min_dy
    ) external {
        address tokenIn = i == 0 ? WETH : USDT;
        address tokenOut = j == 0 ? WETH : USDT;

        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), dx);
        IERC20(tokenIn).safeApprove(CURVE_POOL, dx);

        ICurve(CURVE_POOL).exchange_underlying(i, j, dx, min_dy);

        uint256 dy = IERC20(tokenOut).balanceOf(address(this));
        IERC20(tokenOut).safeTransfer(msg.sender, dy);
    }

    function swapZyber(
        uint256 amount0Out,
        uint256 amount1Out,
        bytes calldata data
    ) external {
        address tokenIn = amount0Out == 0 ? USDT : WETH;
        uint256 amountIn = amount0Out == 0 ? amount1Out : amount0Out;

        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);
        IERC20(tokenIn).safeApprove(ZYBER_POOL, amountIn);
        IZyber(ZYBER_POOL).swap(amount0Out, amount1Out, msg.sender, data);

        address tokenOut = amount0Out == 0 ? WETH : USDT;
        uint256 amountOut = IERC20(tokenOut).balanceOf(address(this));
        IERC20(tokenOut).safeTransfer(msg.sender, amountOut);
    }
}
