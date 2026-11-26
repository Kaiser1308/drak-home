#pragma once

#include <cstddef>
#include <string_view>

namespace drak {

// Maximum length of a contract v1 identifier (command_id, device_id, ...).
inline constexpr std::size_t kMaxIdLength = 128;

// Returns true if s matches the contract v1 identifier pattern
// ^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$ (see contracts/v1/schemas/command.json).
bool is_valid_id(std::string_view s) noexcept;

}  // namespace drak
