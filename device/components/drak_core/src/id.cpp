#include "drak/id.hpp"

namespace drak {
namespace {

constexpr bool is_alnum(char c) noexcept {
    return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9');
}

}  // namespace

bool is_valid_id(std::string_view s) noexcept {
    if (s.empty() || s.size() > kMaxIdLength || !is_alnum(s.front())) {
        return false;
    }
    for (char c : s.substr(1)) {
        if (!is_alnum(c) && c != '_' && c != '.' && c != '-') {
            return false;
        }
    }
    return true;
}

}  // namespace drak
