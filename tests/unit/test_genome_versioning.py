"""Tests for genome schema versioning."""

import pytest
from cards_evolve.genome.versioning import SchemaVersion, validate_schema_version
from cards_evolve.genome.schema import GameGenome


def test_current_schema_version() -> None:
    """Test current schema version is 1.0."""
    assert SchemaVersion.CURRENT == "1.0"


def test_validate_compatible_version() -> None:
    """Test compatible schema versions pass validation."""
    genome = GameGenome(schema_version="1.0", genome_id="test", generation=0)

    # Should not raise
    validate_schema_version(genome)


def test_validate_incompatible_version_raises() -> None:
    """Test incompatible schema versions raise error."""
    genome = GameGenome(schema_version="2.0", genome_id="test", generation=0)

    with pytest.raises(ValueError, match="Incompatible schema version"):
        validate_schema_version(genome)
